#include<stdio.h>


int g=7;


int main(void){
    float f=1;
    int* p=&g;
    printf("%.2f %p %d %p\n",f,p,*p,p+1);
    //g--既是一个逻辑表达式，又是一个赋值语句
    if(g--){//判断g是否为0
        printf("g is not 0\n");
    }else{
        printf("g is 0\n");
    }
}