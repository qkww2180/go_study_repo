#coding=utf-8
__author__='orisun'
 
import numpy as np
import time
 
LEN=10000000
arr=[]
brr=[]
for i in range(LEN):
    value=i%10
    arr.append(value)
    brr.append(value)
array1=np.array(arr)
array2=np.array(brr)
 
begin=time.time()*1000
prod=0
for i in range(LEN):
    prod+=arr[i]*brr[i]
end=time.time()*1000
print("蛮力计算内积 %.2e\t\t用时%d毫秒" % (prod,end-begin))
 
begin=time.time()*1000
prod=np.dot(array1,array2)
end=time.time()*1000
print("使用numpy计算内积 %.2e\t用时%d毫秒" % (prod,end-begin))

# python basic/fast/dp.py