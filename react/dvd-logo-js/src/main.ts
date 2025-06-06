import { setupCanvas } from './canvas'
import './style.css'

document.querySelector<HTMLDivElement>('#app')!.innerHTML = `
  <div>
    <canvas id="dvdCanvas"></canvas>
  </div>
`

const canvas = document.querySelector<HTMLCanvasElement>('#dvdCanvas')
if (!canvas) {
  throw new Error("no canvas")
}

setupCanvas(canvas)
