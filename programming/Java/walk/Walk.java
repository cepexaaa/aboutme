package info.kgeorgiy.ja.kubesh.walk;

import java.io.*;
import java.nio.charset.StandardCharsets;
import java.nio.file.*;
import java.security.MessageDigest;
import java.security.NoSuchAlgorithmException;
import java.util.Arrays;

public class Walk {
    private static final String HASH_ALGORITHM = "SHA-256";
    private static final String ZERO_HASH = "0000000000000000";
    private static final int BUFFER_SIZE = 8192;

    public static void main(String[] args) {
        if (!isArgumentsInvalid(args)) {
            System.err.println("Invalid arguments");
            return;
        }
        String inputFile = args[0];
        String outputFile = args[1];
        try {
            Path inputFilePath = Paths.get(inputFile);
            if (!inputFilePath.toFile().exists()) {
                System.err.println("Input file does not exist: " + inputFile);
                return;
            }
            Path outputFilePath = Paths.get(outputFile);
            if (outputFilePath.getParent() != null) {
                try {
                    Files.createDirectories(outputFilePath.getParent());
                } catch (IOException e) {
                    System.err.println("Failed to create output directory: " + outputFilePath);
                }
            }
            walkFiles(inputFilePath, outputFilePath);
        } catch (InvalidPathException e) {
            System.err.println("Invalid input file path: " + inputFile);
        } catch (NullPointerException e) {
            System.err.println("Null pointer exception: " + e);
        }

    }

    private static void walkFiles(Path inputFilePath, Path outputFilePath) {
        try (BufferedReader reader = Files.newBufferedReader(inputFilePath, StandardCharsets.UTF_8)) {
            try (BufferedWriter writer = Files.newBufferedWriter(outputFilePath, StandardCharsets.UTF_8)) {
                String line;
                while ((line = reader.readLine()) != null) {
                    if (line.isEmpty()) {
                        continue;
                    }
                    try {
                        Path path = Path.of(line);
                        if (Files.isRegularFile(path)) {
                            processFile(path, writer);
                        } else {
                            System.err.println("File doesn't exist: " + line);
                            writer.write(ZERO_HASH + " " + line);
                            writer.newLine();
                        }
                    } catch (IOException e) {
                        writer.write(ZERO_HASH + " " + line);
                        writer.newLine();
                        System.err.println("Error processing file: " + line);
                    } catch (InvalidPathException e) {
                        writer.write(ZERO_HASH + " " + line);
                        writer.newLine();
                        System.err.println("Invalid input file path: " + line);
                    }
                }
            } catch (IOException e) {
                System.err.println("Error by write " + e.getMessage());
            }
        } catch (IOException e) {
            System.err.println("During working with files an error occurred (error occurred during read): " + e.getMessage());
        }
    }

    private static void processFile(Path file, BufferedWriter writer) throws IOException {
        String hash;
        try {
            hash = calculateSHA256(file);
        } catch (IOException e) {
            hash = ZERO_HASH;
        }
        writer.write(hash + " " + file);
        writer.newLine();
    }

    private static String calculateSHA256(Path file) throws IOException {
        try (InputStream inputStream = Files.newInputStream(file)) {
            MessageDigest digest = MessageDigest.getInstance(HASH_ALGORITHM);
            byte[] buffer = new byte[BUFFER_SIZE];
            int bytesRead;
            while ((bytesRead = inputStream.read(buffer)) != -1) {
                digest.update(buffer, 0, bytesRead);
            }
            byte[] hashBytes = digest.digest();
            byte[] last64Bits = Arrays.copyOfRange(hashBytes, 0, 8);
            StringBuilder hexString = new StringBuilder();
            for (byte b : last64Bits) {
                hexString.append(String.format("%02x", b));
            }
            return hexString.toString();
        } catch (IOException e) {
            throw new IOException("During calculating the hash amount error occurred", e);
        } catch (NoSuchAlgorithmException e) {
            throw new NoSuchFileException("No such algorithm");
        }
    }

    private static boolean isArgumentsInvalid(String[] args) {
        return args != null && args.length == 2 && args[0] != null && args[1] != null;
    }
}
