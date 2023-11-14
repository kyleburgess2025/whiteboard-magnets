import { useEffect, useState } from "react";
import AddWordTile from "../components/AddWordTile";
import WordTile from "../components/WordTile";
import useWebSocket from "react-use-websocket";

const WS_URL = "ws://127.0.0.1:8080/ws";

export interface WordTileInformation {
  xValue: number;
  yValue: number;
  word: string;
  id: string;
}

const isGetEvent = (message: any): boolean => {
  return message.type === "get";
};

const isAddEvent = (message: any): boolean => {
  return message.type === "add";
};

function Whiteboard() {
  const [words, setWords] = useState<WordTileInformation[]>([]);
  const [loading, isLoading] = useState<boolean>(true);
  const { sendJsonMessage, readyState, lastJsonMessage } = useWebSocket(
    WS_URL,
    {
      onOpen: () => {
        console.log("WebSocket connection established.");
      },
      share: true,
      filter: () => true,
      retryOnError: true,
      shouldReconnect: (closeEvent) => true,
    }
  );

  useEffect(() => {
    if (readyState === 1) {
      isLoading(false);
      sendJsonMessage({
        type: "get",
      });
    }
  }, [readyState]);

  useEffect(() => {
    type LastJsonMessage = {
      type: string;
      words: WordTileInformation[];
    };
    type AddJsonMessage = {
      type: string;
      word: WordTileInformation;
    };

    if (
      lastJsonMessage &&
      isGetEvent(lastJsonMessage) &&
      (lastJsonMessage as LastJsonMessage).words
    ) {
      const lastJsonMessageTyped = lastJsonMessage as LastJsonMessage;
      console.log(lastJsonMessageTyped.words);
      setWords(lastJsonMessageTyped.words);
    } else if (
      lastJsonMessage &&
      isAddEvent(lastJsonMessage) &&
      (lastJsonMessage as AddJsonMessage).word
    ) {
      const lastJsonMessageTyped = lastJsonMessage as AddJsonMessage;
      setWords([...words, lastJsonMessageTyped.word]);
    }
  }, [lastJsonMessage]);

  const addWord = (newWord: WordTileInformation) => {
    setWords([...words, newWord]);
    sendJsonMessage({
      type: "add",
      word: newWord,
    });
  };

  return (
    <div>
      {words.map((word) => (
        <WordTile
          word={word.word}
          xValue={word.xValue}
          yValue={word.yValue}
          id={word.id}
          sendJsonMessage={sendJsonMessage}
          readyState={readyState}
          lastJsonMessage={lastJsonMessage}
          key={word.id}
        />
      ))}
      <AddWordTile addWord={addWord} />
    </div>
  );
}

export default Whiteboard;
