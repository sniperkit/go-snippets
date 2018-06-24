#import <Foundation/Foundation.h>

int main(void)
{
  int i, j;
  for (i = 1; i <= 12; i++) {
    for (j = 1; j <= 12; j++)
      GSPrintf(stdout, @"%4d", i * j);
    GSPrintf(stdout, @"\n");
  }
  return 0;
}
