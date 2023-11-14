import { useState } from "react";
import { WordTileInformation } from "../pages/Whiteboard";
import useWindowDimensions from "./UseWindowDimensions";
import { v4 as uuidv4 } from "uuid";

interface AddWordTileProps {
  addWord: (value: WordTileInformation) => void;
}

// TODO: Make x and y random pixel values, or center of screen

function AddWordTile({ addWord }: AddWordTileProps) {
  const [word, setWord] = useState<string>("");
  const { height, width } = useWindowDimensions();

  const onAddWord = () => {
    addWord({
      word,
      xValue: Math.floor(Math.random() * (width - 5)),
      yValue: Math.floor(Math.random() * (height - 5)),
      id: uuidv4(),
    });
    setWord("");
  };

  return (
    <div>
      <label htmlFor="word">Enter a word.</label>
      <br />
      <input
        type="text"
        name="word"
        value={word}
        onChange={(e) => setWord(e.target.value)}
      />
      <button onClick={onAddWord}>Add word</button>
    </div>
  );
}

export default AddWordTile;
