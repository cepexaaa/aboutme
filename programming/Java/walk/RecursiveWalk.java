package info.kgeorgiy.ja.kubesh.walk;

import java.io.*;
import java.nio.charset.StandardCharsets;
import java.nio.file.DirectoryStream;
import java.nio.file.Files;
import java.nio.file.Path;
import java.nio.file.Paths;
import java.security.MessageDigest;
import java.util.Arrays;

public class RecursiveWalk {
    public static void main(String[] args) {
        if (args.length != 2) {
            System.err.println("Two input files are expected. Current count of arguments - " + args.length);
        }

        String inputFile = args[0];
        String outputFile = args[1];

        try (BufferedReader reader = Files.newBufferedReader(Paths.get(inputFile), StandardCharsets.UTF_8);
             BufferedWriter writer = Files.newBufferedWriter(Paths.get(outputFile), StandardCharsets.UTF_8)) {

            String line;
            while ((line = reader.readLine()) != null) {
                line = line.trim();
                if (line.isEmpty()) {
                    continue;
                }
                Path path = Paths.get(line);
                if (Files.exists(path)) {
                    if (Files.isDirectory(path)) {
                        processDirectory(path, writer);
                    } else {
                        processFile(path, writer);
                    }
                } else {
                    System.err.println("File or directory doesn't exist: " + line);
                    writer.write("0000000000000000 " + line);
                    writer.newLine();
                }
            }
        } catch (IOException e) {
            System.err.println("During working with files an error occurred: " + e.getMessage());
        }
    }

    private static void processDirectory(Path directory, BufferedWriter writer) throws IOException {
        try (DirectoryStream<Path> stream = Files.newDirectoryStream(directory)) {
            for (Path entry : stream) {
                if (Files.isDirectory(entry)) {
                    processDirectory(entry, writer);
                } else {
                    processFile(entry, writer);
                }
            }
        } catch (IOException e) {
            writer.write("0000000000000000 " + directory);
            writer.newLine();
        }
    }

    private static void processFile(Path file, BufferedWriter writer) throws IOException {
        String hash;
        try {
            hash = calculateSHA256(file);
        } catch (IOException e) {
            hash = "0000000000000000";
        }
        writer.write(hash + " " + file);
        writer.newLine();
    }

    private static String calculateSHA256(Path file) throws IOException {
        try (InputStream inputStream = Files.newInputStream(file)) {
            MessageDigest digest = MessageDigest.getInstance("SHA-256");
            byte[] buffer = new byte[8192];
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
        } catch (Exception e) {
            throw new IOException("During calculating the hash amount error occurred", e);
        }
    }
}

/**

 package info.kgeorgiy.ja.kubesh.walk;

 import java.io.*;
 import java.nio.charset.StandardCharsets;
 import java.nio.file.Files;
 import java.nio.file.Path;
 import java.nio.file.Paths;
 import java.security.MessageDigest;
 import java.util.Arrays;

 public class Walk {
 public static void main(String[] args) {
 if (args.length != 2) {
 System.err.println("Two input files are expected. Current count of arguments - " + args.length);
 }

 String inputFile = args[0];
 String outputFile = args[1];

 try (BufferedReader reader = Files.newBufferedReader(Paths.get(inputFile), StandardCharsets.UTF_8);
 BufferedWriter writer = Files.newBufferedWriter(Paths.get(outputFile), StandardCharsets.UTF_8)) {

 String line;
 while ((line = reader.readLine()) != null) {
 line = line.trim();
 if (line.isEmpty()) {
 continue;
 }
 Path path = Paths.get(line);
 if (Files.isRegularFile(path)) {
 processFile(path, writer);
 } else {
 System.err.println("File doesn't exist: " + line);
 writer.write("0000000000000000 " + line);
 writer.newLine();
 }
 }
 } catch (IOException e) {
 System.err.println("During working with files an error occurred: " + e.getMessage());
 }
 }

 private static void processFile(Path file, BufferedWriter writer) throws IOException {
 String hash;
 try {
 hash = calculateSHA256(file);
 } catch (IOException e) {
 hash = "0000000000000000";
 }
 writer.write(hash + " " + file);
 writer.newLine();
 }

 private static String calculateSHA256(Path file) throws IOException {
 try (InputStream inputStream = Files.newInputStream(file)) {
 MessageDigest digest = MessageDigest.getInstance("SHA-256");
 byte[] buffer = new byte[8192];
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
 } catch (Exception e) {
 throw new IOException("During calculating the hash amount error occurred", e);
 }
 }
 }

 */
