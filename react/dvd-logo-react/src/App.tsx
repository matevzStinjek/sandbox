import { useEffect, useRef, useState } from 'react'
import viteLogo from '/vite.svg'
import './App.css'

function App() {
  const [position, setPosition] = useState({ x: 0, y: 10 })

  const delta = useRef({ dx: 50, dy: 50 })

  const updateFrame = () => {
    setPosition((prevPosition) => {
      let newDx = delta.current.dx;
      let newDy = delta.current.dy;

      const nextX = prevPosition.x + newDx;
      const nextY = prevPosition.y + newDy;

      if (nextX > window.innerWidth || nextX < 0) {
        newDx = -newDx;
      }
      if (nextY > window.innerHeight || nextY < 0) {
        newDy = -newDy;
      }
      delta.current = { dx: newDx, dy: newDy };

      return {
        x: prevPosition.x + newDx,
        y: prevPosition.y + newDy,
      };
    })
  }

  useEffect(() => {
    const intervalId = setInterval(updateFrame, 100)
    return () => clearInterval(intervalId)
  }, [])

  return (
    <>
      <img src={viteLogo} className="logo" alt="Vite logo"  style={{ transform: `translate(${position.x}px, ${position.y}px)` }} />
    </>
  )
}

export default App
