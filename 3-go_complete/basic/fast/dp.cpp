#include <iostream>
#include <ctime>
#include <Windows.h>
 
 
float inner_product(const float* x, const float* y, const long & len){
    float prod = 0.0f;
    long i;
    for (i=0;i<len;i++){
        prod+=x[i]*y[i];
    }
    return prod;
}
 
int main(){
    SetConsoleOutputCP(CP_UTF8); // 设置控制台输出为 UTF-8 编码

    const int len2 = 10000000;
    float *crr=new float[len2];
    float *drr=new float[len2];
    for (int i=0;i<len2;i++){
        int value=i%10;
        crr[i]=value;
        drr[i]=value;
    }
 
    float prod;
    clock_t begin;
    clock_t end;
    begin=GetTickCount(); //从操作系统启动经过的毫秒数
    prod=inner_product(crr,drr,len2);
    end=GetTickCount(); //从操作系统启动经过的毫秒数
    std::cout<<"蛮力计算内积 "<<prod<<"\t用时"<<(double)(end-begin)<<"毫秒"<<std::endl;
}

// g++ -O2 basic/fast/dp.cpp -o basic/fast/dp    不开启编译优化，耗时多一倍
// ./basic/fast/dp