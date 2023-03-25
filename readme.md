<div align="center">
<p>
    <img width="150" src="logos/ArLogo.png">
</p>
<h1>Argon v3</h1>
</div>

ARGON 3 is a math-driven programming language designed to make code easy to read and write. It's not meant to be fast, as it's interpreted. This specification should be used as a guideline, and is subject to change for later versions. Later updates for Argon 3 should be backwards compatible (where possible) with code designed for older versions of the interpreter.

## 📚 Features
   - Easy to read and write: Argon 3 is designed with clarity of code in mind, making it easier for you and others to read and write code.
   - All numbers are stored as rational numbers, preventing precision errors.
   - Math-driven: Designed for mathematical computations, Argon 3 uses techniques and rules set in maths. It's designed to be easy for mathematicians to write and understand algorithms in.
   - Interpreted: Argon 3 is an interpreted language, so you don't need to compile your code before running it.
   - Cross-platform: Argon 3 can be run on any platform that has an interpreter for it.
   - Lightweight: The Argon 3 interpreter is small and doesn't require a lot of system resources to run.

## 💻 Installation
As of now, Argon 3 does not have an installer. Feel free to clone this repo and run the `build` file for your plateform. the build will be found in `bin/argon(.exe)`.

## 📖 Usage
To use Argon 3, you can create a file with the .ar extension and write your code in it. Then, you can run your code using the interpreter. For example, if you have a file called example.ar, you can run it using the following command:

```
argon example.ar
```

## 🔍 Specification

For a detailed specification of the Argon 3 language, please refer to [spec.md](spec.md).

## 🚀 Example Code

Here's an example of how to define a function in Argon 3:

```javascript
f(x) = x^2 + 2*x + 1
term.log('f(10) =', f(10))
```

This code defines a function f(x) that calculates x^2 + 2*x + 1. It then calls the function with an argument of 10 and logs the result to the console.

Please note that this example is subject to change as the specification is in beta and may be updated frequently.

## Licence

MIT
