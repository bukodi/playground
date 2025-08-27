package main

import (
	"bytes"
	"context"
	"crypto/ed25519"
	"crypto/rand"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"golang.org/x/crypto/ssh"
)

func main() {
	// Configuration: replace with your own user, password, and authorized public key as needed.
	const (
		listenAddr      = ":2222"
		allowedUser     = "demo"
		allowedPassword = "s3cret" // set to empty "" to disable password logins
	)
	// Optional: set an authorized_keys-style line for public key auth
	// e.g., `ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIE... comment`
	authorizedPubKeyLine := os.Getenv("AUTHORIZED_PUBKEY") // or hardcode for testing

	// Host key: generate an in-memory Ed25519 key
	_, hostPriv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		log.Fatalf("failed to generate host key: %v", err)
	}
	hostSigner, err := ssh.NewSignerFromSigner(hostPriv)
	if err != nil {
		log.Fatalf("failed to create host signer: %v", err)
	}

	// Build server config
	cfg := &ssh.ServerConfig{
		ServerVersion: "SSH-2.0-GoMiniSSH",
	}
	cfg.AddHostKey(hostSigner)

	// Public key auth
	var authorizedKey ssh.PublicKey
	if strings.TrimSpace(authorizedPubKeyLine) != "" {
		parsed, _, _, _, e := ssh.ParseAuthorizedKey([]byte(authorizedPubKeyLine))
		if e != nil {
			log.Fatalf("failed to parse AUTHORIZED_PUBKEY: %v", e)
		}
		authorizedKey = parsed
		cfg.PublicKeyCallback = func(meta ssh.ConnMetadata, key ssh.PublicKey) (*ssh.Permissions, error) {
			if meta.User() != allowedUser {
				return nil, fmt.Errorf("unknown user")
			}
			if bytes.Equal(key.Marshal(), authorizedKey.Marshal()) {
				return &ssh.Permissions{}, nil
			}
			return nil, fmt.Errorf("unauthorized key")
		}
	}

	// Password auth (optional)
	if allowedPassword != "" {
		cfg.PasswordCallback = func(meta ssh.ConnMetadata, pass []byte) (*ssh.Permissions, error) {
			if meta.User() == allowedUser && string(pass) == allowedPassword {
				return &ssh.Permissions{}, nil
			}
			return nil, errors.New("invalid credentials")
		}
	}

	// If neither callback is set, no one can log in
	if cfg.PublicKeyCallback == nil && cfg.PasswordCallback == nil {
		log.Fatal("no authentication method enabled; set AUTHORIZED_PUBKEY and/or allowedPassword")
	}

	ln, err := net.Listen("tcp", listenAddr)
	if err != nil {
		log.Fatalf("listen error: %v", err)
	}
	log.Printf("SSH server listening on %s (user=%s)", listenAddr, allowedUser)

	// Graceful shutdown
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	go func() {
		<-ctx.Done()
		_ = ln.Close()
	}()

	for {
		nConn, err := ln.Accept()
		if err != nil {
			if ctx.Err() != nil {
				log.Println("shutting down...")
				return
			}
			log.Printf("accept error: %v", err)
			continue
		}
		go handleConn(nConn, cfg)
	}
}

func handleConn(nConn net.Conn, cfg *ssh.ServerConfig) {
	defer nConn.Close()

	sshConn, chans, reqs, err := ssh.NewServerConn(nConn, cfg)
	if err != nil {
		log.Printf("handshake failed: %v", err)
		return
	}
	defer sshConn.Close()
	log.Printf("new ssh connection from %s (%s)", sshConn.RemoteAddr(), sshConn.ClientVersion())

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
			log.Printf("could not accept channel: %v", err)
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
