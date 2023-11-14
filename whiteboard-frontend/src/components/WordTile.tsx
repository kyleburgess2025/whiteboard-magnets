import { useState, useRef, useEffect } from "react";
import { ReadyState } from "react-use-websocket";

interface WordTileProps {
  word: string;
  xValue: number;
  yValue: number;
  id: string;
  sendJsonMessage: (message: any) => void;
  readyState: number;
  lastJsonMessage: any;
}

const isMoveEvent = (message: any): boolean => {
  return message.type === "move";
};

function WordTile(wordTile: WordTileProps) {
  const [isDragging, setDragging] = useState(false);
  const block = useRef<HTMLDivElement>(null);
  const frameID = useRef(0);
  const lastX = useRef(0);
  const lastY = useRef(0);
  const dragX = useRef(0);
  const dragY = useRef(0);
  const { sendJsonMessage, readyState, lastJsonMessage } = wordTile;

  useEffect(() => {
    if (
      lastJsonMessage &&
      isMoveEvent(lastJsonMessage) &&
      lastJsonMessage.word.id &&
      lastJsonMessage.word.id === wordTile.id
    ) {
      dragX.current = lastJsonMessage.word.deltaX;
      dragY.current = lastJsonMessage.word.deltaY;
      cancelAnimationFrame(frameID.current);
      frameID.current = requestAnimationFrame(() => {
        if (block.current === null) {
          return;
        }
        block.current.style.transform = `translate3d(${dragX.current}px, ${dragY.current}px, 0)`;
      });
    }
  }, [lastJsonMessage]);

  const handleMove = (e: MouseEvent) => {
    if (!isDragging) {
      return;
    }

    const deltaX = lastX.current - e.pageX;
    const deltaY = lastY.current - e.pageY;
    lastX.current = e.pageX;
    lastY.current = e.pageY;
    dragX.current -= deltaX;
    dragY.current -= deltaY;
    sendJsonMessage({
      type: "move",
      word: {
        word: wordTile.word,
        id: wordTile.id,
        xValue: block.current?.getBoundingClientRect().left,
        yValue: block.current?.getBoundingClientRect().top,
        deltaX: dragX.current,
        deltaY: dragY.current,
      },
    });

    cancelAnimationFrame(frameID.current);
    frameID.current = requestAnimationFrame(() => {
      if (block.current === null) {
        return;
      }
      block.current.style.transform = `translate3d(${dragX.current}px, ${dragY.current}px, 0)`;
    });
  };

  const handleMouseDown = (e: React.MouseEvent<HTMLElement>) => {
    lastX.current = e.pageX;
    lastY.current = e.pageY;
    setDragging(true);
  };

  const handleMouseUp = () => {
    setDragging(false);
  };

  useEffect(() => {
    document.addEventListener("mousemove", handleMove);
    document.addEventListener("mouseup", handleMouseUp);

    return () => {
      document.removeEventListener("mousemove", handleMove);
      document.removeEventListener("mouseup", handleMouseUp);
    };
  }, [isDragging]);

  // Make it so broadcast to all users except the one who moved it

  return (
    <div>
      <div
        ref={block}
        onMouseDown={handleMouseDown}
        style={{
          position: "absolute",
          left: wordTile.xValue,
          top: wordTile.yValue,
        }}
      >
        {wordTile.word}
      </div>
    </div>
  );
}

export default WordTile;
