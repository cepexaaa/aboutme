package info.kgeorgiy.ja.kubesh.hello;

import info.kgeorgiy.java.advanced.hello.NewHelloClient;

import java.io.IOException;
import java.net.*;
import java.nio.charset.StandardCharsets;
import java.util.List;
import java.util.concurrent.ExecutorService;
import java.util.concurrent.Executors;
import java.util.concurrent.TimeUnit;
import java.util.stream.IntStream;

/**
 * Simple example of client path in connection.
 * Send request to server on UDP.
 * Implements {@link NewHelloClient}
 * {@inheritDoc}
 */
public class HelloUDPClient implements NewHelloClient {
    /**
     * {@inheritDoc}
     */
    @Override
    public void newRun(List<Request> requests, int threads) {
        ExecutorService executor = Executors.newFixedThreadPool(threads);
        IntStream.range(0, threads).forEach(thread -> {
            executor.submit(() -> {
                try (DatagramSocket socket = new DatagramSocket()) {
                    //note -- лучше завести константу типа RECEIVE_TIMEOUT
                    socket.setSoTimeout(100);

                    for (Request request : requests) {
                        String message = request.template().replace("$", String.valueOf(thread + 1));
                        processRequest(socket, request.host(), request.port(), message);
                    }
                } catch (SocketException e) {
                    throw new RuntimeException("Socket creation failed", e);
                }
            });
        });
        close(executor);
    }

    private static void close(ExecutorService executor) {
        //note -- just executor.close()
        executor.shutdown();
        try {
            if (!executor.awaitTermination(1, TimeUnit.MINUTES)) {
                executor.shutdownNow();
            }
        } catch (InterruptedException e) {
            executor.shutdownNow();
            Thread.currentThread().interrupt();
        }
    }


    private void processRequest(DatagramSocket socket, String host, int port, String message) {
        while (true) {
            try {
                InetAddress address = InetAddress.getByName(host);
                byte[] sendData = message.getBytes(StandardCharsets.UTF_8);
                DatagramPacket packet = new DatagramPacket(sendData, sendData.length, address, port);
                socket.send(packet);

                byte[] receiveData = new byte[socket.getReceiveBufferSize()];
                DatagramPacket responsePacket = new DatagramPacket(receiveData, receiveData.length);
                socket.receive(responsePacket);

                String response = new String(responsePacket.getData(), responsePacket.getOffset(), responsePacket.getLength(), StandardCharsets.UTF_8);

                if (response.equals("Hello, " + message)) {
                    return;
                }
            } catch (SocketTimeoutException e) {
                // continue
            } catch (IOException e) {
                throw new RuntimeException("Communication error", e);
            }
        }
    }

    /**
     * Input arguments to set connection.
     * @param args to settings
     */
    public static void main(String[] args) {
        if (args == null || args.length != 5) {
            System.err.println("Usage: HelloUDPClient <host> <port> <prefix> <threadCount> <requestCount>");
            return;
        }

        try {
            String host = args[0];
            int port = Integer.parseInt(args[1]);
            String prefix = args[2];
            int threads = Integer.parseInt(args[3]);
            int requests = Integer.parseInt(args[4]);
            new HelloUDPClient().run(host, port, prefix, requests, threads);
        } catch (NumberFormatException e) {
            System.err.println("Invalid argument: " + e.getMessage());
        }
    }
}
