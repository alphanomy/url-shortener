import { useState } from "react";
import axios, { AxiosResponse } from "axios";
import "./App.css";

function App() {
  const [url, setUrl] = useState("");
  const [shortenedUrl, setShortenedUrl] = useState("");
  const [isCopied, setIsCopied] = useState(false);

  const handleSubmit = async () => {
    try {
      const body = { url: url };

      const response: AxiosResponse = await axios.post("/shorten", body);

      const responseData = response.data;
      setShortenedUrl(`http://127.0.0.1:8000/${responseData.shortened_url}`);
      setIsCopied(false);
    } catch (error) {
      console.error(error);
    }
  };

  const handleCopyShortenedUrl = () => {
    navigator.clipboard.writeText(shortenedUrl);
    setIsCopied(true);
  };

  return (
    <>
      <h1>Create Short URL</h1>
      <div className="container">
        <input
          type="text"
          placeholder="Enter Your URL"
          onChange={(e) => setUrl(e.target.value)}
        />
        <button onClick={handleSubmit} disabled={url ? false : true}>
          Shorten
        </button>
      </div>
      {shortenedUrl && (
        <div className="container">
          <a href={shortenedUrl} target="_blank">
            {shortenedUrl}
          </a>
          <button onClick={handleCopyShortenedUrl}>
            {isCopied ? "Copied" : "Copy"}
          </button>
        </div>
      )}
    </>
  );
}

export default App;
