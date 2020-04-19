/*************************************************************************
    > File Name: io.cpp
    > Author: lixutong
    > Created Time: 2020-04-19 14:07:59
    > Description: 从标准输入进行读取，并写入标准输出
 ************************************************************************/

#include "../include/apue.h"
#include <iostream>
#define BUFFSIZE 4096

int main()
{
    int n;
    char buf[BUFFSIZE];
    while ((n = read(STDIN_FILENO, buf, BUFFSIZE)) > 0)
        if (write(STDOUT_FILENO, buf, n) != n)
            std::cout << "write error!" << std::endl;
    if (n < 0)
        std::cout << "read error!\n" << std::endl;
    exit(0);
}

