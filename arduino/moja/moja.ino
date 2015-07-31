/*
 * Project moja
 *
 * PIN assignment
 *
 * A0(PF0) COL1(R) 25,28,31,34,37,40,43,46
 * A1(PF1) COL2(G) 25,28,31,34,37,40,43,46
 * A2(PF2) COL3(B) 25,28,31,34,37,40,43,46
 * A3(PF3) COL4(R) 26,29,32,35,38,41,44,47
 * A4(PF4) COL5(G) 26,29,32,35,38,41,44,47
 * A5(PF5) COL6(B) 26,29,32,35,38,41,44,47
 * A6(PF6) COL7(R) 27,30,33,36,39,42,45,48
 * A7(PK0) COL9(B) 27,30,33,36,39,42,45,48
 *
 * A8(PF7)
 * 22(PA0) ROW1 1,2,3
 * 23(PA1) ROW2 4,5,6
 * 24(PA2) ROW3 7,8,9
 * 25(PA3) ROW4 10,11,12
 * 26(PA4) ROW5 13,14,15
 * 27(PA5) ROW6 16,17,18
 * 28(PA6) ROW7 19,20,21
 * 29(PA7) ROW8 22,23,24
 * 30(PC7) ROW16 46,47,48
 * 31(PC6) ROW15 43,44,45
 * 32(PC5) ROW14 40,41,42
 * 33(PC4) ROW13 37,38,39
 * 34(PC3) ROW12 34,35,36
 * 35(PC2) ROW11 31,32,33
 * 36(PC1) ROW10 28,29,30
 * 37(PC0) ROW9 25,26,27
 * 38(PD7)
 * 39(PG2)
 * 40(PG1)
 * 41(PG0) COL7(R) 3,6,9,12,15,18,21,24
 * 42(PL7) COL8(G) 3,6,9,12,15,18,21,24
 * 43(PL6) COL9(B) 3,6,9,12,15,18,21,24
 * 44(PL5) COL4(R) 2,5,8,11,14,17,20,23
 * 45(PL4) COL5(G) 2,5,8,11,14,17,20,23
 * 46(PL3) COL6(B) 2,5,8,11,14,17,20,23
 * 47(PL2) COL1(R) 1,4,7,10,13,16,19,22
 * 48(PL1) COL2(G) 1,4,7,10,13,16,19,22
 * 49(PL0) COL3(B) 1,4,7,10,13,16,19,22
 * 50(PB3)
 * 51(PB2)
 * 52(PB1)
 * 53(PB0) 
 */

static HardwareSerial *pSerial = &Serial;  // console for test and degug
//static HardwareSerial *pSerial = &Serial1; // for RasPi

#define BRIGHTNESS_CTRL
//#define COL_ACTIVE_HIGH

#ifdef COL_ACTIVE_HIGH
  #define SET_COL_A(a) \
    PORTL = a; \
    PORTG = (a >> 8) & 0x1
  #define SET_COL_B(b) \
    PORTF = b; \
    PORTK = (b >> 8) & 0x1
  #define CLEAR_RED_A \
    PORTL &= ~0x24; \
    PORTG &= ~0x1
  #define CLEAR_RED_B \
    PORTF &= ~0x24; \
    PORTK &= ~0x1
#else
  #define SET_COL_A(a) \
    PORTL = ~a; \
    PORTG = ~((a >> 8) & 0x1)
  #define SET_COL_B(b) \
    PORTF = ~b; \
    PORTK = ~((b >> 8) & 0x1)
  #define CLEAR_RED_A \
    PORTL |= 0x24; \
    PORTG |= 0x1
  #define CLEAR_RED_B \
    PORTF |= 0x24; \
    PORTK |= 0x1
#endif

/*
 * [7] blinking bit
 * [3..6] Reserved
 * [2] color bit R
 * [1] color bit G
 * [0] color bit B
 * 1 001 blue
 * 2 010 green
 * 3 011 cyan
 * 4 100 red
 * 5 101 purple
 * 6 110 yellow
 * 7 111 white
 */
static uint8_t led[47];

void setup()
{
  pSerial->begin(115200);

  // IO port
  DDRA = 0xFF;  // PORTA for ROW1-8
  DDRL = 0xFF;  // PORTL for COL1-8
  DDRG = 0x01;  // PORTG for COL9
  
  DDRC = 0xFF;  // PORTC for ROW9-16
  DDRF = 0xFF;  // PORTF for COL1-8
  DDRK = 0x01;  // PORTF for COL9
    
  TCCR1A = 0;
  TCCR1B = 0;
  // 16MHz / 64 / 10 = 25kHz
  bitSet(TCCR1B, WGM12);  // CTC mode
  bitSet(TCCR1B, CS10);   // 1/64
  bitSet(TCCR1B, CS11);   // 1/64
  OCR1A = 10 - 1;
  bitSet(TIMSK1, OCIE1A); // set mask
  sei();
  
  pinMode(13,OUTPUT);
}

#define CMDBUF_LEN 128

char cmdbuf[CMDBUF_LEN];
uint8_t cmdbuf_idx = 0;

uint8_t read_line(HardwareSerial *serial)
{
  char c;
  
  while(serial->available()){
    c = serial->read();
    if(c == '\r') goto eol;
    if(c <= ' ') continue;
    cmdbuf[cmdbuf_idx] = c;
    cmdbuf_idx++;
    if(cmdbuf_idx == CMDBUF_LEN - 1) goto eol;
  }
  return 0;  
eol:
  cmdbuf[cmdbuf_idx] = '\0';
  uint8_t r = cmdbuf_idx;
  cmdbuf_idx = 0;
  return r;
}

int8_t next_int(char **_p)
{
  char buf[3], *p = *_p;
  
  if(p == NULL || p >= cmdbuf + CMDBUF_LEN) return -1;
  buf[0] = p[0];
  if(p[1] == ',' || p[1] == '\0'){
    buf[1] = '\0';
    *_p += 2;
  }else{
    buf[1] = p[1];
    if(p[2] == ',' || p[2] == '\0'){
      buf[2] = '\0';
      *_p += 3;
    }else{
      return -1;  // error
    }
  }
  // EOL?
  // is previous NULL?
  if(*(*_p - 1) == '\0'){
    *_p = NULL;
  }
  return atoi(buf);
}

void cmd_parse()
{
  char *p = cmdbuf;
  int8_t n, pref_id, led_stat;
  uint8_t blink_bit = 0;
  
  // type
  switch(next_int(&p)){
  case 0: break;
  case 1: blink_bit = 0x80; break;
  default: // error
    pSerial->println("type error");
    return;
  }
  
  do{
    pref_id = next_int(&p);
    led_stat = next_int(&p);
    if(pref_id >= 0 && led_stat >= 0){
      if(pref_id < 1 || 47 < pref_id){
        pSerial->print("pref_id error:");
        pSerial->println(pref_id);
        return;
      }
      led[pref_id - 1] = blink_bit + (uint8_t)led_stat;
      pSerial->print(pref_id);
      pSerial->print(":");
      pSerial->println(led_stat);
    }
  }while(pref_id >= 0 && led_stat >= 0);
}

void loop()
{
  if(read_line(pSerial)){
    pSerial->println(cmdbuf);
    cmd_parse();
  }
}

static uint16_t row_bit = 1;
static uint8_t row_num = 0;

static uint8_t set_col_a(boolean blk)
{
  // offset
  uint8_t os = 3 * row_num, i, n = 0;
  uint16_t a = 0;
  
  // 1 - 24
  if(!(led[os] & 0x80) || (led[os] & 0x80 && blk)){
    a |= led[os] & 0x7;  // [0..2]
  }
  if(!(led[os + 1] & 0x80) || (led[os + 1] & 0x80 && blk)){
    a |= (led[os + 1] & 0x7) << 3; // [3..5]
  }
  if(!(led[os + 2] & 0x80) || (led[os + 2] & 0x80 && blk)){
    a |= (uint16_t)(led[os + 2] & 0x7) << 6; // [6..8]
  }
  SET_COL_A(a);
  
  for(i = 0;i < 9; i++){
    if(a & 0x01) n++;
    a >>= 1;
  }
  return n;
}
  
static uint8_t set_col_b(boolean blk)
{
  uint8_t os = 3 * row_num, i, n = 0;
  uint16_t b = 0;

  // 25 -48
  if(!(led[24 + os] & 0x80) || (led[24 + os] & 0x80 && blk)){
    b |= led[24 + os] & 0x7; // [0..2]
  }
  if(!(led[24 + os + 1] & 0x80) || (led[24 + os + 1] & 0x80 && blk)){
    b |= (led[24 + os + 1] & 0x7) << 3; // [3..5]
  }
  if(!(led[24 + os + 2] & 0x80) || (led[24 + os + 2] & 0x80 && blk)){
    b |= (uint16_t)(led[24 + os + 2] & 0x7) << 6; // [6..8]
  }
  SET_COL_B(b);
  
  for(i = 0;i < 9; i++){
    if(b & 0x01) n++;
    b >>= 1;
  }
  return n;
}

ISR(TIMER1_COMPA_vect)
{
  static uint16_t cnt = 0;
  static boolean blk = HIGH;
  static uint8_t an = 0, bn = 0, da = 15, db = 15, ta = 0, tb = 0;
  
  if(cnt % 47 == 0){
    // all off
    SET_COL_A(0);    // LED 1 - 24
    SET_COL_B(0);    // LED 25 - 48
  
    // set row
    PORTA = (uint8_t)row_bit;
    PORTC = (uint8_t)row_bit;
    an = set_col_a(blk);
    bn = set_col_b(blk);
    ta = 12 + an * 4;
    tb = 12 + bn * 4;
    
#ifdef BRIGHTNESS_CTRL
    da = ta / 4;
    if(da > 15) da = 15;
    db = tb / 4;
    if(db > 15) db = 15;
#endif
    // next row  
    row_bit <<= 1;
  
    row_num ++;
    if(row_bit == 0x100){
      row_bit = 1;
      row_num = 0;
    }
#ifdef BRIGHTNESS_CTRL
  }else{
    if(cnt % 47 == ta){
      SET_COL_A(0);
    }
    if(cnt % 47 == tb){
      SET_COL_B(0);
    }
    if(cnt % 47 == da){ // red dimmer
      CLEAR_RED_A;
    }
    if(cnt % 47 == db){ // red dimmer
      CLEAR_RED_B;
    }
  }
#else
  }else if(cnt % 47 == 15){
    CLEAR_RED_A;
    CLEAR_RED_B;
  }
#endif

  if(cnt++ == 12500){
    blk = !blk;
    digitalWrite(13, blk);
    cnt = 0;
  }
}
