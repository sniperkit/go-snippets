#import <Foundation/Foundation.h>

int main(void)
{
  int i;
  for (i = 1; i < 100; i += 2) {
    GSPrintf(stdout, @"%i\n", i);
  }
  return 0;
}
