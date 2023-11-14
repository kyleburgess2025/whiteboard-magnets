import { Link } from "react-router-dom";

function Home() {
  return (
    <div>
      <h1>Welcome! Join a fridge...</h1>
      <Link to="/whiteboard">
        <button>Click me!</button>
      </Link>
    </div>
  );
}

export default Home;
