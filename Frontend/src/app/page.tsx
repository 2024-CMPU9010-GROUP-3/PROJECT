import axios from "axios";
import HomePage from "./components/home/HomePage";

export default function Home() {
  const fetchFakeData = async () => {
    const data = await axios.get("https://jsonplaceholder.typicode.com/posts");
    console.log(data);
  };
  fetchFakeData();
  return (
    <div>
      <HomePage />
    </div>
  );
}
