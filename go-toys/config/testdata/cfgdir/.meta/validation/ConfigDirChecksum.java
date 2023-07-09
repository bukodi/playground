import java.math.BigInteger;
import java.nio.ByteBuffer;
import java.nio.charset.StandardCharsets;
import java.nio.file.*;
import java.security.MessageDigest;
import java.util.Arrays;
import java.util.stream.Collectors;
import java.util.stream.Stream;

class ConfigDirChecksum {

    static Path startDir = Paths.get("./../..");
    static byte[] checksum = new byte[32];

    public static void main(String[] args) throws Exception {
        processDir(startDir);

        System.out.printf("Checksum : %s\n", toHex(checksum));
    }

    static void processDir(Path dir) throws Exception {
        for (Path item : Files.list(dir).collect(Collectors.toList())) {
            if (item.getFileName().toString().startsWith("."))
                continue;

            if (Files.isDirectory(item)) {
                processDir(item);
                continue;
            }

            MessageDigest hash = MessageDigest.getInstance("SHA-256");
            String relativeToStart = startDir.relativize(item).toString();
            hash.update(relativeToStart.getBytes(StandardCharsets.UTF_8));
            hash.update(Files.readAllBytes(item));
            byte[] itemHash = hash.digest();

            System.out.printf("%s %s\n", toHex(itemHash), relativeToStart);

            xorItemHash(itemHash);
        }
    }

    static void xorItemHash(byte[] itemHash) {
        for (int i = 0; i < 32; i++) {
            checksum[i] = (byte) (checksum[i] ^ itemHash[i]);
        }
    }

    static String toHex(byte[] bytes) {
        ByteBuffer buffer = ByteBuffer.wrap(bytes);
        return Stream.generate(() -> buffer.get()).
                limit(bytes.length).
                map(b -> String.format("%02x", b)).
                collect(Collectors.joining());
    }

}