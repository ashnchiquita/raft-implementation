export const get = async (key: string) => {
    const response = await fetch(`/app/${key}`);
    return response.json();
  };
  
  export const set = async (key: string, value: string) => {
    const response = await fetch('/app', {
      method: 'PUT',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({ key, value }),
    });
    return response.json();
  };
  
  export const del = async (key: string) => {
    const response = await fetch(`/app/${key}`, {
      method: 'DELETE',
    });
    return response.json();
  };
  
  export const append = async (key: string, value: string) => {
    const response = await fetch('/app', {
      method: 'PATCH',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({ key, value }),
    });
    return response.json();
  };
  
  export const getAll = async () => {
    const response = await fetch('/app');
    return response.json();
  };
  
  export const delAll = async () => {
    const response = await fetch('/app', {
      method: 'DELETE',
    });
    return response.json();
  };
  
  export const strlen = async (key: string) => {
    const response = await fetch(`/app/${key}/strlen`);
    return response.json();
  };
  