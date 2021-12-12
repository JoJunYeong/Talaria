/*
 * Author: Jo Jun Yeong
 * Example For sending data to SerialPort
 */

#define BAUD 9600
#define DELAY_TIME 1000

int pin2 = 2;
int pin3 = 3;
int pin4 = 4;
int pin5 = 5;
int pin6 = 6;
int pin7 = 7;
int pin8 = 8;
int pin9 = 9;
int pin10 = 10;
int pin11 = 11;
int pin12 = 12;
int pin13 = 13;
int pin22 = 22;
int pin23 = 23;


int velocity = 0;
int send_count =0;
int reset_complete=0;
byte leds = 0;
unsigned int count=0;

const char* headingway = "";

void setup() {
  
  pinMode(pin2, INPUT);     // 방향Data 들어오는 첫번째 pin
  pinMode(pin3, INPUT);     // 방향Data 들어오는 두번째 pin
  pinMode(pin4, INPUT);     // 속도Data 들어오는 첫번째 pin
  pinMode(pin5, INPUT);     // 속도 pin
  pinMode(pin6, INPUT);     // 속도 pin
  pinMode(pin7, INPUT);     // 속도 pin
  pinMode(pin8, INPUT);     // 속도 pin
  pinMode(pin9, INPUT);     // 속도 pin
  pinMode(pin10, INPUT);    // 속도 pin
  pinMode(pin11, INPUT);    // 속도Data 들어오는 마지막 pin
  pinMode(pin12, INPUT);    // 속도데이터 다 보냈다는 시그널 pin
  
  pinMode(pin13, OUTPUT);    // 
  pinMode(pin22, OUTPUT);    // 
  pinMode(pin23, INPUT);
  Serial.begin(BAUD);


  digitalWrite(pin2, LOW); 
  digitalWrite(pin3, LOW); 
  digitalWrite(pin4, LOW); 
  digitalWrite(pin5, LOW); 
  digitalWrite(pin6, LOW); 
  digitalWrite(pin7, LOW); 
  digitalWrite(pin8, LOW); 
  digitalWrite(pin9, LOW); 
  digitalWrite(pin10, LOW); 
  digitalWrite(pin11, LOW); 
  digitalWrite(pin12, LOW); 
  digitalWrite(pin13, HIGH); 
  digitalWrite(pin23, LOW);

  
}


void sequence_reset(){
  if(digitalRead(pin23) == HIGH){

    digitalWrite(pin22, HIGH);
    digitalWrite(pin13, HIGH);
    send_count=0;
    reset_complete=1;
    
  }
  
}





void loop() {

 sequence_reset();
  if(reset_complete==1){
    reset_complete=0;
    goto re;
  }
 re:
  velocity = 0;
 if(send_count==0 && (digitalRead(pin13) == LOW))   // 버그 방지
 {
  digitalWrite(pin13, HIGH); 
  sequence_reset();
  if(reset_complete==1){
    reset_complete=0;
    goto re;
  }
 }
  if(digitalRead(pin12) == LOW)
  {
    sequence_reset();
    if(reset_complete==1){
      reset_complete=0;
      goto re;
    }
    count++;
    if(count%60000==0&&count!=0){       // count==60000 일때 low 출력을 시킨다. 아두이노가 생각보다 연산이 빠르기 때문에 이러한 조치를 내렸다.
 //   Serial.print("digitalRead(pin12) == LOW    count = ");
 //   Serial.println(count);
    }
  }


 if(digitalRead(pin12) == HIGH)  // 12번핀에 데이터가 들어오면 속도함수가 다 전송된거임
  {
    digitalWrite(pin22, LOW); 
    if(send_count==0)
    {
      
      digitalWrite(pin13, LOW);       // 그리고 Green LED를 킨다.
      if(digitalRead(pin12) == HIGH)  // 
        {
          while(digitalRead(pin12) == HIGH)    // digitalRead(pin12) == HIGH 가 거짓이 될때까지 while 루프
          {
            sequence_reset();
            if(reset_complete==1){
              reset_complete=0;
              goto re;
            }
            /*
            Serial.print("digitalRead(pin12) = " );
            Serial.println(digitalRead(pin12));
            
            Serial.print("send_count = ");
            Serial.println(send_count);
            
            Serial.println("loop");
            */
            
            digitalRead(pin12);
            if(digitalRead(pin12)==LOW)
              break;
            }
            sequence_reset();
            if(reset_complete==1){
              reset_complete=0;
              goto re;
            }
        }
      //  delay(5000);
      send_count=1;
 //     Serial.println("send_count 0->1");
    }
    else if(send_count==1)
    {
  //     Serial.print("send_count = ");
 //      Serial.println(send_count);
            
 //      Serial.println("enter else if(send_count==1)");
      /////////////////////// 방향수신코드 시작 /////////////////////////
       // 어차피 디지털 핀 코드를 읽어오는거니까 라즈베리파이에서 모든 정보를 정제해서 인코드 하면 이녀석이 아두이노가 디코드 한다.
       if(digitalRead(pin3) == LOW && digitalRead(pin2) == LOW)
       {
        // 전진방향 전송코드작성 필요
        headingway = "f";
       }
       else if(digitalRead(pin3) == LOW && digitalRead(pin2) == HIGH) 
       {
        // 후진방향 전송코드 작성 필요 
        headingway = "b";
       }
       else if(digitalRead(pin3) == HIGH && digitalRead(pin2) == LOW) 
       {
        // 왼쪽방향 전송코드 작성 필요 
        headingway = "l";
       }
       else if(digitalRead(pin3) == HIGH && digitalRead(pin2) == HIGH) 
       {
        // 오른쪽방향 전송코드 작성 필요 
        headingway = "r";
       }
       /////////////////////// 방향수신코드 끝 /////////////////////////
  
  
     /////////////////////// 속도계산코드 시작 /////////////////////////
/*
      Serial.print(digitalRead(pin11));
      Serial.print(digitalRead(pin10));
      Serial.print(digitalRead(pin9));
      Serial.print(digitalRead(pin8));
      Serial.print(digitalRead(pin7));
      Serial.print(digitalRead(pin6));
      Serial.print(digitalRead(pin5));
      Serial.println(digitalRead(pin4));
  */
   
       if(digitalRead(pin4) == HIGH)    //pin4는 2의 0제곱을 의미한다. HIGH로 왔으면 활성화 시킨것이므로 velocity에 넣는다.
        velocity = velocity + 1;
       if(digitalRead(pin5) == HIGH)    //pin5는 2의 1제곱을 의미한다. HIGH로 왔으면 활성화 시킨것이므로 velocity에 넣는다.
        velocity = velocity + 2;
       if(digitalRead(pin6) == HIGH)    //pin6는 2의 2제곱을 의미한다. HIGH로 왔으면 활성화 시킨것이므로 velocity에 넣는다.
        velocity = velocity + 4;
       if(digitalRead(pin7) == HIGH)    //pin7는 2의 3제곱을 의미한다. HIGH로 왔으면 활성화 시킨것이므로 velocity에 넣는다.
        velocity = velocity + 8;
       if(digitalRead(pin8) == HIGH)    //pin8는 2의 4제곱을 의미한다. HIGH로 왔으면 활성화 시킨것이므로 velocity에 넣는다.
        velocity = velocity + 16;
       if(digitalRead(pin9) == HIGH)    //pin9는 2의 5제곱을 의미한다. HIGH로 왔으면 활성화 시킨것이므로 velocity에 넣는다.
        velocity = velocity + 32;
       if(digitalRead(pin10) == HIGH)    //pin10는 2의 6제곱을 의미한다. HIGH로 왔으면 활성화 시킨것이므로 velocity에 넣는다.
        velocity = velocity + 64;
       if(digitalRead(pin11) == HIGH)    //pin11는 2의 7제곱을 의미한다. HIGH로 왔으면 활성화 시킨것이므로 velocity에 넣는다.
        velocity = velocity + 128;
  
  
             
        Serial.print(headingway);     // 시리얼 포트로 타깃 컴퓨터에 움직이는방향 전송
        Serial.print("@");     // 방향데이터와 속도데이터의 구분을 알리는 기호 전송
        Serial.print(velocity);     // 시리얼 포트로 타깃컴퓨터에 속도 전송
        Serial.print("!");   // 데이터 전송이 끝났음을 알리는 기호 전송
 //       Serial.println("#");    
        digitalWrite(pin22, HIGH);       // 그리고 Green LED를 킨다.
       // delay(DELAY_TIME);
       send_count=0;

  //    Serial.println("send_count 1->0");
      
       while(digitalRead(pin12) == HIGH)    // digitalRead(pin12) == HIGH 가 거짓이 될때까지 while 루프
       {
        digitalRead(pin12);
          if(digitalRead(pin12) == HIGH)
          {
            sequence_reset();
            if(reset_complete==1){
              reset_complete=0;
              goto re;
            }
            digitalRead(pin12);
            count++;
 //           if(count%60000==0&&count!=0){       // count==60000 일때 low 출력을 시킨다. 아두이노가 생각보다 연산이 빠르기 때문에 이러한 조치를 내렸다.
 //           Serial.print("digitalRead(pin12) ==");
 ///           Serial.print(digitalRead(pin12));
  //          Serial.print("          digitalRead(pin12) == HIGH    count = ");
  //          Serial.println(count);
  //          }
          }
      }
      sequence_reset();
      if(reset_complete==1){
        reset_complete=0;
        goto re;
      }
      digitalWrite(pin13, HIGH);       // 그리고 Green LED를 킨다.
 //     Serial.println("pin13 low->high");
    }
        sequence_reset();
        if(reset_complete==1){
          reset_complete=0;
          goto re;
        }
 }
   /////////////////////// 속도계산코드 끝 /////////////////////////

sequence_reset();
if(reset_complete==1){
    reset_complete=0;
    goto re;
  }
//
//digitalWrite(pin13, HIGH);
//delay(5000);
//digitalWrite(pin13, LOW);
//delay(5000);


//////////////////// 이 지점에 이르러서는 속도와 방향까지 전부 다 계산되어있다. //////////////////////////

//  if(velocity!=0)   // 이 말은 loop앞에서 velocity를 0으로 초기화를 해놓았으니 0이 아니라는 것은 속도값이 들어왔다는 것을 의미한다. 그렇기때문에 그 속도와 방향을 Serialport로 전송해야 한다.
//  {
//    Serial.println(velocity);     // 시리얼 포트로 속도 전송
//    Serial.println(headingway);     // 시리얼 포트로 진행방향 전송
//    delay(DELAY_TIME);
//  }























 
}
