// Allow importing text files with Vite's ?raw query, e.g. import txt from './file.txt?raw'
declare module '*?raw' {
  const value: string;
  export default value;
}

// Keep a more specific pattern for plain .txt raw imports as well.
declare module '*.txt?raw' {
  const value: string;
  export default value;
}

// For imports without ?raw (plain text modules, rare) - treat as string too
declare module '*.txt' {
  const value: string;
  export default value;
}
