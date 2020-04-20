/*************************************************************************
    > File Name: pid.cpp
    > Author: lixutong
    > Created Time: 2020-04-19 15:39:44
    > Description: print pid of this process
 ************************************************************************/

#include "apue.h"
#include <iostream>
using namespace std;

int main()
{   
    cout << "process id " << static_cast<long>(getpid()) << endl;
    exit(0);
}

