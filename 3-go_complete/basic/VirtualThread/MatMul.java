package basic.VirtualThread;

import java.util.Random;
import java.util.concurrent.ExecutorService;
import java.util.concurrent.Executors;
import java.util.concurrent.TimeUnit;

public class MatMul{
    // 任务本身
    private static Runnable runnable = () -> {
    };

    public static void awaitTerminationAfterShutdown(ExecutorService threadPool) {
        threadPool.shutdown(); 
        try {
            if (!threadPool.awaitTermination(60, TimeUnit.SECONDS)) {
                threadPool.shutdownNow();
            }
        } catch (InterruptedException ex) {
            threadPool.shutdownNow();
            Thread.currentThread().interrupt();
        }
    }

    public static void main(String[] args) {
        // 热身，第一次不算
        int p=Integer.parseInt(args[0]);//p个并发线程
        System.out.println("concurrency "+p);
        ExecutorService executorService = Executors.newScheduledThreadPool(10);//定长线程池
        for (int i = 0; i < 10000; i++) {  //JIT warm up 1万次 
            executorService.submit(runnable);
        }
        try{TimeUnit.SECONDS.sleep(1);}catch(InterruptedException e){};   //等热身任务线束

        // 第二次才算。复用ExecutorService
        long startTime = System.currentTimeMillis();//开始计时
        for (int i = 0; i < p; i++) {
            executorService.submit(runnable);
        }
        awaitTerminationAfterShutdown(executorService);//等所有任务结束
        long endTime = System.currentTimeMillis();//结束计时
        System.out.println("use time "+(endTime-startTime)+" ms");
    }
} 

// javac -encoding UTF-8 basic/VirtualThread/MatMul.java 
// java basic/VirtualThread/MatMul 10000