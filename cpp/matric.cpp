#include <stdio.h> 

int main() {
  char matric[8];
  int q1, q2, q3, q4;
  float sum, avg;
  float totalAvg = 0;
  int count = 0;

  while (1) {
    printf("Enter student %d matric number :", count + 1);
    scanf("%s", matric);

    // break loop
    if (matric[0] == '0') {
    	break;
    }

	printf("Enter quiz 1:");
	scanf("%d", &q1);
	printf("Enter quiz 2:");
	scanf("%d", &q2);
	printf("Enter quiz 3:");
	scanf("%d", &q3);
	printf("Enter quiz 4:");
	scanf("%d", &q4);

	// calculate sum
	sum = q1 + q2 + q3 + q4;
	printf("Sum is %.0f\n", sum);

	// calculate avg
	avg = sum / 4;
	totalAvg += avg;
	
	count++;
  }

  totalAvg /= count;
  printf("Number of students: %d\n", count);
  printf("Average marks for this class is : %.2f\n", totalAvg);
  // system("pause");
  return 0;
}