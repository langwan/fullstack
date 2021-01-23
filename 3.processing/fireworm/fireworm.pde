
Star[] stars;

int screenWidth = 800;
int screenHeight = 400;
int starTotal =  int(random(400, 1000));

class Star {
  float radius;
  float light;
  PVector trackPos;
  float trackRaduis;
  float angle;
  float frequency;
}

void setup() {
  size(800 , 400);
  noStroke();
  rectMode(CENTER);
  stars = new Star[starTotal];
  for (int i = 0; i < starTotal; i++) {
    stars[i] = new Star();
    stars[i].trackPos = new PVector(random(screenWidth), random(screenHeight));
    stars[i].trackRaduis = random(10, 120);
    stars[i].radius = random(3);
    stars[i].light = random(100, 300);
    stars[i].angle = random(360);
    stars[i].frequency = random(0.01, 0.03);
  }
}

void draw() {
  background(0);
  for (int i = 0; i < starTotal; i++) {
    Star star = stars[i];
    float light = star.light * random(0.5, 1);
    fill(255, light);
    
    float x = star.trackPos.x + cos(radians(star.angle))  * star.trackRaduis;
    float y = star.trackPos.y + sin(radians(star.angle))  * star.trackRaduis;
    
    
    ellipse(x, y, star.radius, star.radius);
    star.angle -= star.frequency;
  } 
}
