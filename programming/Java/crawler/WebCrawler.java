package info.kgeorgiy.ja.kubesh.crawler;

import info.kgeorgiy.java.advanced.crawler.*;

import java.io.IOException;
import java.net.MalformedURLException;
import java.util.*;
import java.util.concurrent.*;

/**
 * Walk links in site.
 * Visit all links in page site. In each page visit their links.
 * Implements interface {@link Crawler}.
 * {@inheritDoc}
 *
 * @see Crawler
 */
public class WebCrawler implements NewCrawler {
    private static final int DEFAULT_CRAWL_VALUES = 10;
    private final Downloader downloader;
    private final ExecutorService downloadersPool;
    private final ExecutorService extractorsPool;

    /**
     * Constructor for setting crawl parameters.
     *
     * @param downloader  allows you to download pages and extract links from them
     * @param downloaders maximum number of simultaneously loaded pages
     * @param extractors  the maximum number of pages from which links can be extracted at the same time
     * @param perHost     the maximum number of pages loaded simultaneously from a single host. To determine the host, use the getHost method of the URLUtils class from the tests.
     */
    public WebCrawler(Downloader downloader, int downloaders, int extractors, int perHost) {
        this.downloader = downloader;
        this.downloadersPool = Executors.newFixedThreadPool(downloaders);
        this.extractorsPool = Executors.newFixedThreadPool(extractors);
    }

    /**
     * {@inheritDoc}
     */
    @Override
    public Result download(String url, int depth, List<String> excludes) {
        Set<String> visitedUrls = ConcurrentHashMap.newKeySet();
        Map<String, IOException> errors = new ConcurrentHashMap<>();
        Set<String> downloaded = ConcurrentHashMap.newKeySet();

        Phaser phaser = new Phaser(1);
        Queue<Task> queue = new ConcurrentLinkedQueue<>();
        queue.add(new Task(url, depth));

        while (true) {
            int tasksThisRound = queue.size();
            if (tasksThisRound == 0) {
                break;
            }
            processFS(excludes, phaser, tasksThisRound, queue, visitedUrls, errors, downloaded);

            phaser.arriveAndAwaitAdvance();
        }

        return new Result(new ArrayList<>(downloaded), errors);
    }

    private void processFS(List<String> excludes, Phaser phaser, int tasksThisRound, Queue<Task> queue, Set<String> visitedUrls, Map<String, IOException> errors, Set<String> downloaded) {
        phaser.bulkRegister(tasksThisRound);
        for (int i = 0; i < tasksThisRound; i++) {
            Task task = queue.poll();
            String currentUrl = task.url;
            int currentDepth = task.depth;

            if (currentDepth <= 0 || !isValidURL(currentUrl, excludes) || !visitedUrls.add(currentUrl)) {
                phaser.arriveAndDeregister();
                continue;
            }
            processDownload(task.url, task.depth, phaser, errors, downloaded, queue);
        }
    }

    private static class Task {
        final String url;
        final int depth;
        Task(String url, int depth) {
            this.url = url;
            this.depth = depth;
        }
    }

    private void processLinks(Document document, String url, int depth, Phaser phaser, Map<String, IOException> errors, Queue<Task> queue) {
        executeTask(() -> {
            List<String> links = document.extractLinks();
            links.forEach(link -> queue.add(new Task(link, depth - 1)));
            return null;
        }, url, phaser, errors, extractorsPool);
    }

    private void processDownload(String url, int depth, Phaser phaser, Map<String, IOException> errors, Set<String> downloaded, Queue<Task> queue) {
        executeTask(() -> {
            Document document = downloader.download(url);
            downloaded.add(url);

            if (depth > 1) {
                phaser.register();
                processLinks(document, url, depth, phaser, errors, queue);
            }
            return null;
        }, url, phaser, errors, downloadersPool);
    }

    private <T> void executeTask(Callable<T> task, String url, Phaser phaser, Map<String, IOException> errors, ExecutorService executor) {
        executor.submit(() -> {
            try {
                return task.call();
            } catch (IOException e) {
                errors.put(url, e);
            } finally {
                phaser.arriveAndDeregister();
            }
            return null;
        });
    }

    /**
     * {@inheritDoc}
     */
    @Override
    public Result download(String url, int depth) {
        return download(url, depth, new ArrayList<>());
    }

    private boolean isValidURL(String urlStr, List<String> excludes) {
        try {
            String host = URLUtils.getHost(urlStr);
            for (String exclude : excludes) {
                if (host.contains(exclude)) {
                    return false;
                }
            }
        } catch (MalformedURLException e) {
            return false;
        }
        return true;
    }

    /**
     * {@inheritDoc}
     */
    @Override
    public void close() {
        shutdownPool(downloadersPool, "Downloaders pool");
        shutdownPool(extractorsPool, "Extractors pool");
    }

    private void shutdownPool(ExecutorService pool, String poolName) {
        pool.shutdown();
        try {
            if (!pool.awaitTermination(1, TimeUnit.MINUTES)) {
                pool.shutdownNow();
                System.err.println(poolName + " didn't terminate");
            }
        } catch (InterruptedException e) {
            pool.shutdownNow();
            Thread.currentThread().interrupt();
        }
    }

    private static void startWebCrawler(String url, int depth, int downloads, int extractors, int perHost) {
        try (Crawler crawler = new WebCrawler(new CachingDownloader(1.0), downloads, extractors, perHost)) {
            Result result = crawler.download(url, depth);
            System.out.println("Downloaded " + result.getDownloaded().size() + " pages");
            if (!result.getErrors().isEmpty()) {
                System.out.println("Errors (" + result.getErrors().size() + "):");
                result.getErrors().forEach((u, e) -> System.out.println(u + ": " + e.getMessage()));
            }
        } catch (IOException e) {
            System.err.println("Error creating downloader: " + e.getMessage());
        }
    }

    private static int setNextValue(String[] args, int index) {
        if (args.length > index && args[index].matches("\\d+")) {
            return Integer.parseInt(args[index]);
        }
        return DEFAULT_CRAWL_VALUES;
    }

    /**
     * Start Web crawler from terminal.
     * Accepts 5 parameters to configure the crawl.
     *
     * @param args - downloader pages, number of parallel download, number of pages to extract links and max number of load pages at one time
     */
    public static void main(String[] args) {
        if (args == null || args.length > 5 || args.length < 1) {
            System.err.println("Usage: WebCrawler url [depth [downloads [extractors [perHost]]]]");
            return;
        }

        String url = args[0];
        if (!isCorrectUrl(url)) {
            System.err.println("Invalid URL: " + url);
            return;
        }
        int depth = setNextValue(args, 1);
        int downloaders = setNextValue(args, 2);
        int extractors = setNextValue(args, 3);
        int perHost = setNextValue(args, 4);
        startWebCrawler(url, depth, downloaders, extractors, perHost);

    }

    private static boolean isCorrectUrl(String url) {
        return url.startsWith("http://") || url.startsWith("https://");
    }
}
