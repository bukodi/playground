package main

import (
	"bytes"
	"context"
	"crypto/ed25519"
	"crypto/rand"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net"
	"os"
	"os/exec"
	"strings"
	"syscall"
	"time"

	"golang.org/x/crypto/ssh"
)

type authorizedPubKeyLine string

type SSHServer struct {
	AuthorizedPubKeyLines []authorizedPubKeyLine
	ListenAddr            string
	AllowedUser           string
	AllowedPassword       string
	AllowPQCKex           bool
	AllowNonPQCKex        bool
	ln                    net.Listener
	serveCh               chan struct{}
}

func (s *SSHServer) Start(ctx context.Context) error {
	if s.ListenAddr == "" {
		s.ListenAddr = ":2222"
	}

	// Generate a host key (ephemeral in-memory)
	_, hostPriv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return fmt.Errorf("failed to generate host key: %w", err)
	}
	hostSigner, err := ssh.NewSignerFromSigner(hostPriv)
	if err != nil {
		return fmt.Errorf("failed to create host signer: %w", err)
	}

	kexAlgos := make([]string, 0)
	for _, algo := range ssh.SupportedAlgorithms().KeyExchanges {
		isPqcAlgo := strings.Contains(strings.ToLower(algo), "mlkem")
		if s.AllowPQCKex && isPqcAlgo {
			kexAlgos = append(kexAlgos, algo)
		}
		if s.AllowNonPQCKex && !isPqcAlgo {
			kexAlgos = append(kexAlgos, algo)
		}
	}

	cfg := &ssh.ServerConfig{
		Config: ssh.Config{
			KeyExchanges: kexAlgos,
		},
		ServerVersion: "SSH-2.0-GoMiniSSH",
	}
	cfg.AddHostKey(hostSigner)
	slog.Info("Allowed kex", "kexAlgos", cfg.KeyExchanges)

	// Build authorized keys list
	var authorizedKeys []ssh.PublicKey
	for _, line := range s.AuthorizedPubKeyLines {
		str := strings.TrimSpace(string(line))
		if str == "" {
			continue
		}
		parsed, _, _, _, e := ssh.ParseAuthorizedKey([]byte(str))
		if e != nil {
			return fmt.Errorf("failed to parse authorized pubkey line: %w", e)
		}
		authorizedKeys = append(authorizedKeys, parsed)
	}
	if len(authorizedKeys) > 0 {
		cfg.PublicKeyCallback = func(meta ssh.ConnMetadata, key ssh.PublicKey) (*ssh.Permissions, error) {
			if s.AllowedUser != "" && meta.User() != s.AllowedUser {
				return nil, fmt.Errorf("unknown user")
			}
			for _, k := range authorizedKeys {
				if bytes.Equal(key.Marshal(), k.Marshal()) {
					return &ssh.Permissions{}, nil
				}
			}
			return nil, fmt.Errorf("unauthorized key")
		}
	}

	if s.AllowedPassword != "" {
		cfg.PasswordCallback = func(meta ssh.ConnMetadata, pass []byte) (*ssh.Permissions, error) {
			if (s.AllowedUser == "" || meta.User() == s.AllowedUser) && string(pass) == s.AllowedPassword {
				return &ssh.Permissions{}, nil
			}
			return nil, errors.New("invalid credentials")
		}
	}

	ln, err := net.Listen("tcp", s.ListenAddr)
	if err != nil {
		return fmt.Errorf("listen error: %w", err)
	}
	s.ln = ln
	s.serveCh = make(chan struct{})
	slog.Info("SSH server listening", "addr", s.ListenAddr, "user", s.AllowedUser)

	// Serve in background; respect context cancellation
	go func() {
		defer close(s.serveCh)
		// Stop on context cancellation
		go func() {
			<-ctx.Done()
			_ = ln.Close()
		}()
		for {
			nConn, err := ln.Accept()
			if err != nil {
				// If closed due to Stop() or ctx cancellation, exit serve loop
				if ctx.Err() != nil {
					slog.Info("ssh server shutting down", "reason", "ctx cancelled")
					return
				}
				if ne, ok := err.(net.Error); ok && !ne.Temporary() {
					return
				}
				if strings.Contains(strings.ToLower(err.Error()), "closed") {
					return
				}
				slog.Error("accept error", "err", err)
				continue
			}
			go handleConn(nConn, cfg)
		}
	}()

	return nil
}

func (s *SSHServer) Stop() {
	if s.ln != nil {
		_ = s.ln.Close()
	}
	if s.serveCh != nil {
		<-s.serveCh
	}
}

func main() {
	srv := SSHServer{}
	err := srv.Start(context.Background())
	if err != nil {
		slog.Error("failed to start SSH server", "err", err)
		os.Exit(1)
	}
	time.Sleep(time.Second * 10)
	srv.Stop()
}

func handleConn(nConn net.Conn, cfg *ssh.ServerConfig) {
	defer nConn.Close()

	sshConn, chans, reqs, err := ssh.NewServerConn(nConn, cfg)
	if err != nil {
		slog.Error("handshake failed", "err", err)
		return
	}
	defer sshConn.Close()
	slog.Info("new ssh connection", "remote", sshConn.RemoteAddr().String(), "client_version", string(sshConn.ClientVersion()))

	// Discard global requests
	go ssh.DiscardRequests(reqs)

	// Service channels
	for ch := range chans {
		if ch.ChannelType() != "session" {
			ch.Reject(ssh.UnknownChannelType, "unsupported channel type")
			continue
		}
		channel, requests, err := ch.Accept()
		if err != nil {
			slog.Error("could not accept channel", "err", err)
			continue
		}
		go handleSession(channel, requests)
	}
}

func handleSession(channel ssh.Channel, requests <-chan *ssh.Request) {
	defer channel.Close()

	// Keep some simple session state
	env := map[string]string{}
	var wantShell bool
	var commandToRun string

	for req := range requests {
		switch req.Type {
		case "env":
			// payload: name string, value string
			var kv struct {
				Name  string
				Value string
			}
			_ = ssh.Unmarshal(req.Payload, &kv)
			env[kv.Name] = kv.Value
			req.Reply(true, nil)

		case "pty-req":
			// We don’t implement PTY in this minimal server. Reject but continue.
			req.Reply(false, nil)

		case "shell":
			// No command means user wants an interactive shell
			wantShell = true
			req.Reply(true, nil)
			go runShell(channel, env, "") // basic shell without PTY
			return

		case "exec":
			// payload: command string
			var ex struct {
				Command string
			}
			if err := ssh.Unmarshal(req.Payload, &ex); err != nil {
				req.Reply(false, nil)
				continue
			}
			commandToRun = ex.Command
			req.Reply(true, nil)
			exitCode := runCommand(channel, env, commandToRun)
			sendExitStatus(channel, exitCode)
			return

		case "subsystem":
			// Could implement SFTP here (subsystem "sftp")
			req.Reply(false, nil)

		default:
			req.Reply(false, nil)
		}
	}

	// If we exit loop without shell/exec, just close
	if wantShell && commandToRun == "" {
		_ = sendExitStatus(channel, 0)
	}
}

func runCommand(channel ssh.Channel, env map[string]string, cmd string) int {
	c := exec.Command("/bin/sh", "-lc", cmd)
	c.Stdin = channel
	c.Stdout = channel
	c.Stderr = channel.Stderr()

	// inherit current env + session env
	c.Env = append(os.Environ(), flattenEnv(env)...)

	if err := c.Start(); err != nil {
		io.WriteString(channel, "failed to start command: "+err.Error()+"\n")
		return 127
	}
	err := c.Wait()
	return exitCodeFromErr(err)
}

func runShell(channel ssh.Channel, env map[string]string, shell string) {
	if shell == "" {
		shell = "/bin/sh"
	}
	c := exec.Command(shell)
	c.Stdin = channel
	c.Stdout = channel
	c.Stderr = channel.Stderr()
	c.Env = append(os.Environ(), flattenEnv(env)...)

	if err := c.Start(); err != nil {
		io.WriteString(channel, "failed to start shell: "+err.Error()+"\n")
		sendExitStatus(channel, 127)
		return
	}
	err := c.Wait()
	sendExitStatus(channel, exitCodeFromErr(err))
}

func sendExitStatus(channel ssh.Channel, code int) bool {
	type exitStatus struct {
		Status uint32
	}
	ok, _ := channel.SendRequest("exit-status", false, ssh.Marshal(&exitStatus{Status: uint32(code)}))
	_ = channel.CloseWrite()
	// give client time to read before close
	time.AfterFunc(200*time.Millisecond, func() { _ = channel.Close() })
	return ok
}

func flattenEnv(m map[string]string) []string {
	out := make([]string, 0, len(m))
	for k, v := range m {
		out = append(out, k+"="+v)
	}
	return out
}

func exitCodeFromErr(err error) int {
	if err == nil {
		return 0
	}
	var ee *exec.ExitError
	if errors.As(err, &ee) {
		if status, ok := ee.Sys().(syscall.WaitStatus); ok {
			return status.ExitStatus()
		}
	}
	return 1
}
