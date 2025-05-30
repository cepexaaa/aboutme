package info.kgeorgiy.ja.kubesh.hello;

import info.kgeorgiy.java.advanced.hello.HelloServer;

import java.io.IOException;
import java.net.*;
import java.nio.charset.StandardCharsets;
import java.util.concurrent.*;

/**
 * Simple server to answer on request.
 * Response is "Hello, \<request\>"
 * Implements {@link HelloServer}
 *
 * @see HelloServer
 */
public class HelloUDPServer implements HelloServer {
    private ExecutorService receiverExecutor;
    private ExecutorService workerExecutor;
    private DatagramSocket socket;
    private volatile boolean isRunning;

    /**
     * {@inheritDoc}
     */
    @Override
    public void start(int port, int threads) {
        workerExecutor = Executors.newFixedThreadPool(threads);
        receiverExecutor = Executors.newSingleThreadExecutor();
        isRunning = true;

        try {
            socket = new DatagramSocket(port);
            receiverExecutor.submit(() -> {
                while (!Thread.interrupted() && isRunning) {
                    try {
                        byte[] buffer = new byte[1024];
                        DatagramPacket packet = new DatagramPacket(buffer, buffer.length);
                        socket.receive(packet);

                        workerExecutor.submit(() -> processRequest(packet));
                    } catch (SocketException e) {
                        if (isRunning) {
                            System.err.println("Socket error: " + e.getMessage());
                        }
                    } catch (IOException e) {
                        System.err.println("I/O error: " + e.getMessage());
                    }
                }
            });
        } catch (SocketException e) {
            System.err.println("Could not start server on port " + port + ": " + e.getMessage());
            close();
        }
    }

    private void processRequest(DatagramPacket packet) {
        try {
            String request = new String(packet.getData(), 0, packet.getLength(), StandardCharsets.UTF_8);
            String response = "Hello, " + request;
            byte[] responseData = response.getBytes(StandardCharsets.UTF_8);
            DatagramPacket responsePacket = new DatagramPacket(responseData, responseData.length, packet.getAddress(), packet.getPort());
            socket.send(responsePacket);
        } catch (IOException e) {
            System.err.println("Error processing request: " + e.getMessage());
        }
    }

    /**
     * {@inheritDoc}
     */
    @Override
    public void close() {
        isRunning = false;
        if (socket != null) {
            socket.close();
        }
        if (workerExecutor != null) {
            workerExecutor.shutdownNow();
        }
        if (receiverExecutor != null) {
            receiverExecutor.shutdownNow();
        }
    }

    /**
     * Input arguments to set connection
     * @param args to settings
     */
    public static void main(String[] args) {
        if (args == null || args.length != 2) {
            System.err.println("Usage: HelloUDPServer <port> <threads>");
            return;
        }
        try (HelloUDPServer server = new HelloUDPServer()) {
            int port = Integer.parseInt(args[0]);
            int threads = Integer.parseInt(args[1]);
            server.start(port, threads);
        } catch (NumberFormatException e) {
            System.err.println("Port and threads must be integers: " + e.getMessage());
        }
    }
}
























