class Code {

	public int minFunction(int n1, int n2) {
	   int min;
	   if (n1 > n2)
	      min = n2;
	   else
	      min = n1;

	   return min; 
	}

	public static void printMax( double... numbers) {
	      if (numbers.length == 0) {
		 System.out.println("No argument passed");
		 return;
	      }

	      double result = numbers[0];

	      for (int i = 1; i <  numbers.length; i++){
		 if (numbers[i] >  result){
		      result = numbers[i];
		      System.out.println("The max value is " + result);
		}
	      }
	     
       }

	public void reverse(){
		int num=0;
	      int reversenum =0;
	      System.out.println("Input your number and press enter: ");
	      //This statement will capture the user input
	      Scanner in = new Scanner(System.in);
	      //Captured input would be stored in number num
	      num = in.nextInt();
	      //While Loop: Logic to find out the reverse number
	      while( num != 0 )
	      {
		  reversenum = reversenum * 10;
		  reversenum = reversenum + num%10;
		  num = num/10;
	      }

	      System.out.println("Reverse of input number is: "+reversenum);
	}

	public static boolean isPrime(int num) {
		if (num % 2 == 0){
			return false;
		} 
		for (int i = 3; i * i <= num; i += 2){
		 if (num % i == 0){ return false;}
		}
		   
		return true;
  	} 

	public static void printMoreThan(int a){
		
	    for (int i : new int[]{0, 1, 2, 3, 4, 5, 6, 7, 9}) {
           	 if (i >a){
			System.out.println(i);
		 }
	    }
	
	}

	public static void printTriangle(int a){
		for (int i = 0; i < a; i++)
		{
		    for (int j = a; j > i; j--)
		    {
		        System.out.print(" ");
		    }
		    for (int k = 1; k <= i + 1; k++) {
		        System.out.print(" *");
		    }
		    System.out.print("\n");
		}
	}
}
