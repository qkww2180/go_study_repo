package algorithm

import (
	"bufio"
	"cmp"
	"container/heap"
	"dqq/util/logger"
	"io"
	"os"
	"slices"
	"sort"
	"strconv"
	"strings"
)

func partition[T cmp.Ordered](slice []T, k int, result []T) []T {
	if len(slice) <= k { //递归退出条件
		return slice
	}

	pivot := 0
	i := 1
	j := len(slice) - 1
	for i < j { //一直循环，直到i和j相遇
		//先移动j
		for ; i < j; j-- {
			if slice[j] < slice[pivot] { // slice[0]作为基准
				break
			}
		}
		//后移动i
		for ; i < j; i++ {
			if slice[i] > slice[pivot] {
				break
			}
		}
		//i和j还没有相遇，则交换slice[i]和slice[j]
		if i < j {
			slice[i], slice[j] = slice[j], slice[i]
		}
	}
	//i和j相遇后，如果slice[i] < slice[pivot]，则i和pivot需要交换
	if slice[i] < slice[pivot] {
		//交换slice[pivot]和slice[i]
		slice[pivot], slice[i] = slice[i], slice[pivot]
		pivot = i
	}

	m := len(slice) - pivot
	n := m - 1
	if m == k { //pivot加上它右边的元素刚好有k个
		// return slice[pivot:]
		result = append(result, slice[pivot:]...)
		// result=slice[pivot:]
	} else if n == k { //pivot右边的元素刚好有k个
		result = append(result, slice[pivot+1:]...)
	} else if n > k { //pivot右边的元素多于k个，递归从pivot右边找topK
		result = partition(slice[pivot+1:], k, result)
	} else { //pivot加上它右边的元素少于k个，则pivot及它右边的元素都属于topK，然后再从pivot左边找top  k-(len(slice)-pivot)
		result = append(result, slice[pivot:]...)
		result = partition(slice[:pivot], k-m, result)
	}
	return result
}

func TopKByPartition[T cmp.Ordered](list []T, k int) []T {
	result := make([]T, 0, k)
	brr := slices.Clone(list) //拷贝一份list，别对list造成影响
	return partition(brr, k, result)
}

func TopKByHeap[T cmp.Ordered](list []T, k int) []T {
	if k <= 0 {
		return nil
	}
	if k >= len(list) {
		return list
	}

	heap := NewHeap[T](list[:k]) //用前k个元素构建小根堆
	heap.Build()
	for _, ele := range list[k:] { //依次遍历第k个位置往后的元素
		top, _ := heap.Top()
		if ele > top {
			heap.ReplaceTop(ele) //如果比堆顶大，则替换掉堆顶（堆内部会自行调整）
		}
	}
	return heap.GetAll()
}

// 一个超级大的文件里存着很多ip，一行一个，求出现次数最多的前k个ip。限制：你只有1G的内存
func FindFreqIpFromBigFile(file string, k int) []*Item[int] {
	logger.SetLogLevel(logger.InfoLevel)
	logger.SetLogFile("topk.log")

	const SMALL_FILE_LINE = 10000 //小文件里包含多少个ip，有限的内存只能处理这么小的文件
	fin, err := os.Open(file)
	if err != nil {
		logger.Error("open file %s failed: %s", file, err)
		return nil
	}
	reader := bufio.NewReader(fin)

	//把大文件顺序切分成很多小文件
	var fileCount, lineCount int //最终fileInx的值就是切分成了多少个小文件
	fout, err := os.OpenFile(file+strconv.Itoa(fileCount), os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0o644)
	if err != nil {
		logger.Error("open file %s failed: %s", file+strconv.Itoa(fileCount), err)
		return nil
	}
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err != io.EOF {
				logger.Error("read file %s failed %s", file, err)
			}
			if len(line) > 0 {
				fout.WriteString(line)
			}
			break
		} else {
			fout.WriteString(line)
			lineCount++
			if lineCount >= SMALL_FILE_LINE {
				fout.Close()
				fileCount++
				fout, err = os.OpenFile(file+strconv.Itoa(fileCount), os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0o644)
				if err != nil {
					logger.Error("open file %s failed: %s", file+strconv.Itoa(fileCount), err)
					return nil
				}
				lineCount = 0
			}
		}
	}
	fin.Close()
	fout.Close()

	//对每个小文件按ip从小到大排序，排序后输出到一个文件
	for i := 0; i < fileCount; i++ {
		fin, err = os.Open(file + strconv.Itoa(i))
		if err != nil {
			logger.Error("open file %s failed: %s", file+strconv.Itoa(i), err)
			return nil
		}
		reader := bufio.NewReader(fin)
		lines := make([]string, 0, 10000)
		for {
			line, err := reader.ReadString('\n')
			if err != nil {
				if err != io.EOF {
					logger.Error("read file %s failed %s", file, err)
				}
				if len(line) > 0 {
					lines = append(lines, line)
				}
				break
			} else {
				line = strings.TrimRight(line, "\n")
				lines = append(lines, line)
			}
		}
		fin.Close()
		sort.Slice(lines, func(i, j int) bool {
			return lines[i] < lines[j]
		})
		fout, err = os.OpenFile(file+strconv.Itoa(i)+".sort", os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0o644)
		if err != nil {
			logger.Error("open file %s failed: %s", file+strconv.Itoa(i)+".sort", err)
			return nil
		}
		for _, line := range lines {
			fout.WriteString(line + "\n")
		}
		fout.Close()
	}

	//利用小堆根，对多个有序的小文件进行多路归并排序，输出一个有序的大文件
	fins := make([]*os.File, fileCount)
	readers := make([]*bufio.Reader, fileCount)
	pq := make(PriorityQueue[string], 0, fileCount)
	fileClosed := make([]bool, fileCount)
	//先把每个文件中读出首行，构建初始的小根堆
	for i := 0; i < fileCount; i++ {
		fin, err = os.Open(file + strconv.Itoa(i) + ".sort")
		if err != nil {
			logger.Error("open file %s failed: %s", file+strconv.Itoa(i)+".sort", err)
			return nil
		}
		fins[i] = fin
		reader := bufio.NewReader(fin)
		readers[i] = reader
		if line, err := reader.ReadString('\n'); err == nil {
			pq.Push(&Item[string]{Info: strconv.Itoa(i), Value: strings.TrimRight(line, "\n")})
		}
	}
	fout, err = os.OpenFile(file+".sort", os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0o644)
	if err != nil {
		logger.Error("open file %s failed: %s", file+".sort", err)
		return nil
	}
	heap.Init(&pq) //构建小顶堆
	total := 0
	for pq.Len() > 0 {
		//step1. 把堆顶元素输出至文件，并删除堆顶
		top := heap.Pop(&pq).(*Item[string]) //取出堆顶，并删除，新的堆顶会自动诞生
		fout.WriteString(top.Value + "\n")
		total++
		//step2. 从根堆顶相同的文件里取出下一行，添加到堆中
		idx, _ := strconv.Atoi(top.Info)
		if !fileClosed[idx] { //每读完一个小文件，堆中的元素就会少一个
			line, err := readers[idx].ReadString('\n')
			if err != nil {
				if err != io.EOF {
					logger.Error("read file %s failed %s", file+strconv.Itoa(idx)+".sort", err)
				} else {
					if len(line) > 0 {
						heap.Push(&pq, &Item[string]{Info: strconv.Itoa(idx), Value: line})
					}
				}
				fins[idx].Close()
				fileClosed[idx] = true
			} else {
				heap.Push(&pq, &Item[string]{Info: strconv.Itoa(idx), Value: strings.TrimRight(line, "\n")})
			}
		}
	}
	fout.Close()

	//遍历排序好的大文件，统计每个ip出现的次数，同时用小根堆维护出现次数最多的前k个ip
	prevIp := ""
	prevCount := 0
	queue := make(PriorityQueue[int], 0, k)
	heap.Init(&queue)
	fin, err = os.Open(file + ".sort")
	if err != nil {
		logger.Error("open file %s failed: %s", file+".sort", err)
		return nil
	}
	reader = bufio.NewReader(fin)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				if len(line) > 0 { //最后一行有内容
					if line == prevIp { //老ip
						prevCount++
					} else { //新ip
						prevIp = line
						prevCount = 1
					}
				}
				//把最后一个ip及其计数放入堆中（当然也可能是丢弃）
				if queue.Len() < k {
					heap.Push(&queue, &Item[int]{Info: prevIp, Value: prevCount})
				} else {
					if prevCount > queue[0].Value { //当ip出现的次数大于堆顶时，删除堆顶，把该ip插入堆中
						heap.Pop(&queue)
						heap.Push(&queue, &Item[int]{Info: prevIp, Value: prevCount})
					}
				}
			} else {
				logger.Error("read file %s failed %s", file+".sort", err)
			}
			break
		} else {
			line = strings.TrimRight(line, "\n")
			if line == prevIp {
				prevCount++
			} else {
				if queue.Len() < k {
					heap.Push(&queue, &Item[int]{Info: prevIp, Value: prevCount})
				} else {
					if prevCount > queue[0].Value { //当ip出现的次数大于堆顶时，删除堆顶，把该ip插入堆中
						heap.Pop(&queue)
						heap.Push(&queue, &Item[int]{Info: prevIp, Value: prevCount})
					}
				}
				prevIp = line
				prevCount = 1
			}
		}
	}
	fin.Close()

	rect := make([]*Item[int], 0, k)
	for queue.Len() > 0 {
		top := heap.Pop(&queue).(*Item[int])
		rect = append(rect, top)
	}
	reverse := make([]*Item[int], k) //rect里面是从小到大的，把它逆序过来
	for i := 0; i < k; i++ {
		reverse[i] = rect[k-i-1]
	}
	return reverse
}
