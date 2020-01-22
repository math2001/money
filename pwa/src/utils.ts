export function cast<T>(instance: any, type: { new (...args: any[]): T }): T {
  if (instance instanceof type) return instance;
  console.error({ instance, type: type.name });
  throw new Error("type cast exception");
}

export function assert(condition: boolean, message: string): asserts condition {
  if (condition === false) {
    debugger;
    throw new Error(`Assertion Error: ${message}`);
  }
}
