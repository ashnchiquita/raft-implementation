import { useState } from "react";
import * as apiService from "./apiService";

const EndpointTester = () => {
  const [key, setKey] = useState("");
  const [value, setValue] = useState("");
  const [response, setResponse] = useState(null);

  const handleGet = async () => {
    const res = await apiService.get(key);
    setResponse(res);
  };

  const handleSet = async () => {
    const res = await apiService.set(key, value);
    setResponse(res);
  };

  const handleDelete = async () => {
    const res = await apiService.del(key);
    setResponse(res);
  };

  const handleAppend = async () => {
    const res = await apiService.append(key, value);
    setResponse(res);
  };

  const handleGetAll = async () => {
    const res = await apiService.getAll();
    setResponse(res);
  };

  const handleDelAll = async () => {
    const res = await apiService.delAll();
    setResponse(res);
  };

  const handleStrlen = async () => {
    const res = await apiService.strlen(key);
    setResponse(res);
  };

  return (
    <div>
      <h1>Endpoint Tester</h1>
      <div>
        <input
          type="text"
          placeholder="Key"
          value={key}
          onChange={(e) => setKey(e.target.value)}
        />
        <input
          type="text"
          placeholder="Value"
          value={value}
          onChange={(e) => setValue(e.target.value)}
        />
      </div>
      <div>
        <button onClick={handleGet}>GET</button>
        <button onClick={handleSet}>SET</button>
        <button onClick={handleDelete}>DELETE</button>
        <button onClick={handleAppend}>APPEND</button>
        <button onClick={handleGetAll}>GET ALL</button>
        <button onClick={handleDelAll}>DELETE ALL</button>
        <button onClick={handleStrlen}>STRLEN</button>
      </div>
      <div>
        <h2>Response</h2>
        <pre>{JSON.stringify(response, null, 2)}</pre>
      </div>
    </div>
  );
};

export default EndpointTester;
