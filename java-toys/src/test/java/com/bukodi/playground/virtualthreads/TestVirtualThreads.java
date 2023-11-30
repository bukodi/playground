package com.bukodi.playground.virtualthreads;

import org.junit.Test;

import java.util.concurrent.ExecutorService;

public class TestVirtualThreads {

    @Test
    public void testVirtualThreads() throws InterruptedException {
        System.out.printf("Java version: %s\n", System.getProperty("java.version"));

        try( ExecutorService executorService = java.util.concurrent.Executors.newVirtualThreadPerTaskExecutor() ) {
            ThreadLocal<Integer> currentNum_TL = ThreadLocal.withInitial(() -> 45);

            executorService.execute(() -> {
                currentNum_TL.set(100);
                System.out.printf("Current num: %d\n", currentNum_TL.get());
            });
            executorService.execute(() -> {
                currentNum_TL.set(200);
                System.out.printf("Current num: %d\n", currentNum_TL.get());
            });

            executorService.shutdownNow();
            if( executorService.awaitTermination(5, java.util.concurrent.TimeUnit.SECONDS) ) {
                System.out.print("Finished.\n");
            } else {
                System.out.print("Termination timeout.\n");
            }
        } finally {
            System.out.print("Done.\n");
        }


    }
}
