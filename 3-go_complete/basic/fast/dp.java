package basic.fast;

public class dp{
    private static float InnerProduct(float[] x, float[] y, int len){
        float rect = 0.0f;
        for(int i=0;i<len;i++){
            rect+=x[i]*y[i];
        }
        return rect;
    }

    public static void main(String[] args) {
        final int len=10000000;
        float[] crr = new float[len];
        float[] drr = new float[len];
        for(int i=0;i<len;i++){
            float value = i%10;
            crr[i] = value;
            drr[i] = value;
        }

        for (int i = 0; i < 100; i++) {  //JIT warm up 100次 
            InnerProduct(crr, drr, len);
        }

        long startTime = System.currentTimeMillis();
        float dp = InnerProduct(crr, drr, len);
        long endTime = System.currentTimeMillis();
        System.out.printf("蛮力计算内积 %.0f\t\t用时%d毫秒\n", dp, (int)(endTime-startTime));
    }
}

// javac -encoding UTF-8 basic/fast/dp.java 
// java basic/fast/dp