
Star[] stars;

int screenWidth = 800;
int screenHeight = 400;
int starTotal =  int(random(400, 1000));
class Star {
  float radius;
  float light;
  PVector pos;
}

void setup() {
  size(800 , 400);
  noStroke();
  rectMode(CENTER);
  stars = new Star[starTotal];
  for (int i = 0; i < starTotal; i++) {
    stars[i] = new Star();
    stars[i].pos = new PVector(random(screenWidth), random(screenWidth));
    stars[i].radius = random(3);
    stars[i].light = random(100, 300);
  }
}

void draw() {
  background(0);
  for (int i = 0; i < starTotal; i++) {
    Star star = stars[i];
    float light = star.light * random(0.5, 1);
    fill(255, light);
    println("start", star.pos.x, star.pos.y, star.radius);
    ellipse(star.pos.x, star.pos.y, star.radius, star.radius);
  } 
}
