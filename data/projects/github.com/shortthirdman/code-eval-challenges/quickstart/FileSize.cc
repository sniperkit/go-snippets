/*
File Size  Share on LinkedIn

Description:

Print the size of a file in bytes.

Input sample:

Path to a filename. e.g. 

foo.txt
Output sample:

Print the size of the file in bytes.
e.g.

55
*/

#include <iostream>
#include <sys/types.h>
#include <sys/stat.h>
#include <unistd.h>

using namespace std;

int main( int argc, char** argv ) {
  struct stat filestatus;
  stat( argv[ 1 ], &filestatus );

  cout << filestatus.st_size << endl;

  return 0;
}