#include <NewPing.h>

#define TRIGGER_PIN  4  
#define ECHO_PIN     5  
#define MAX_DISTANCE 100

NewPing sonar(TRIGGER_PIN, ECHO_PIN, MAX_DISTANCE); 

void setup() {
  Serial.begin(9600); 
}

void loop() {
  delay(2000);                     
  unsigned int d = sonar.ping_cm(); 
  Serial.print(d); 
}

