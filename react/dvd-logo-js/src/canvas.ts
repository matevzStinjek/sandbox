export function setupCanvas(canvas: HTMLCanvasElement) {
  canvas.width = 600;
  canvas.height = 400;

  const handleCanvasClick = (event: MouseEvent) => {
    const rect = canvas.getBoundingClientRect();
    const mouseX = event.clientX - rect.left;
    const mouseY = event.clientY - rect.top;
    addLogo(mouseX, mouseY);
  }

  const logos = []

  function addLogo(x: number, y: number) {
    const logoWidth = 80;
    const logoHeight = 40;
    // Center the rect on the click
    const newLogo = new DvdLogo(
      x - logoWidth / 2,
      y - logoHeight / 2,
      logoWidth,
      logoHeight
    );
    logos.push(newLogo);
  }

  canvas.addEventListener("click", handleCanvasClick);
}

const getRandomColor = () => {
  const r = Math.floor(Math.random() * 200 + 55);
  const g = Math.floor(Math.random() * 200 + 55);
  const b = Math.floor(Math.random() * 200 + 55);
  return `rgb(${r},${g},${b})`;
}

class DvdLogo {
  public x: number;
  public y: number;
  public width: number;
  public height: number;
  public dx: number; // Velocity in x
  public dy: number; // Velocity in y
  public color: string;
  public cornerHits: number;

  constructor(
    x: number,
    y: number,
    width?: number, // Optional: if not provided, uses default
    height?: number, // Optional: if not provided, uses default
    initialColor?: string // Optional: if not provided, uses random
  ) {
    this.x = x;
    this.y = y;
    this.width = width || 80; // Default width
    this.height = height || 40; // Default height

    // Random initial direction and speed
    this.dx = (Math.random() > 0.5 ? 1 : -1) * (Math.random() * 1.5 + 1.0); // Speed between 1 and 2.5
    this.dy = (Math.random() > 0.5 ? 1 : -1) * (Math.random() * 1.5 + 1.0);

    this.color = initialColor || getRandomColor();
    this.cornerHits = 0;
  }

  public draw(ctx: CanvasRenderingContext2D): void {
    ctx.fillStyle = this.color;
    ctx.fillRect(this.x, this.y, this.width, this.height);

    // Optional: Display corner hit count
    ctx.fillStyle = "white"; // Text color
    ctx.font = "12px Arial";
    // Position text inside the rectangle, near the top-left
    ctx.fillText(`Corners: ${this.cornerHits}`, this.x + 5, this.y + 15);
  }

  public update(canvasWidth: number, canvasHeight: number): void {
    this.x += this.dx;
    this.y += this.dy;

    let hitVerticalEdge = false;
    let hitHorizontalEdge = false;

    // Collision with left/right walls
    if (this.x + this.width > canvasWidth) {
      this.x = canvasWidth - this.width; // Snap to edge
      this.dx *= -1;
      hitVerticalEdge = true;
    } else if (this.x < 0) {
      this.x = 0; // Snap to edge
      this.dx *= -1;
      hitVerticalEdge = true;
    }

    // Collision with top/bottom walls
    if (this.y + this.height > canvasHeight) {
      this.y = canvasHeight - this.height; // Snap to edge
      this.dy *= -1;
      hitHorizontalEdge = true;
    } else if (this.y < 0) {
      this.y = 0; // Snap to edge
      this.dy *= -1;
      hitHorizontalEdge = true;
    }

    // Corner Hit Detection
    if (hitVerticalEdge && hitHorizontalEdge) {
      console.log("CORNER HIT!");
      this.cornerHits++;
      this.color = getRandomColor(); // Change color on corner hit
    }
  }
}