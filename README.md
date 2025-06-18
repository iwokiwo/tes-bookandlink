# IWO Frontend

Frontend project built with [React](https://reactjs.org/), [Vite](https://vitejs.dev/), [Material UI](https://mui.com/), [Redux Toolkit](https://redux-toolkit.js.org/), [React Hook Form](https://react-hook-form.com/), and [React Router](https://reactrouter.com/).

# Penjelasan
```
Sistem ini tidak hanya memproses antrean (queue) secara efisien, tetapi juga secara dinamis mensimulasikan proses ping ke setiap URL dalam antrean. Jika ping ke URL tertentu gagal, maka sistem secara otomatis akan menandai status antrean tersebut sebagai "failed".
```

## 🔧 Tech Stack

- **Framework**: React 18 + Vite
- **UI Library**: MUI v6 + Emotion
- **Form Validation**: React Hook Form + Yup
- **Routing**: React Router DOM v7
- **State Management**: Redux Toolkit
- **Async State**: React Query
- **Date Handling**: Moment.js
- **API**: Axios
- **Linting**: ESLint + TypeScript ESLint
- **Type Checking**: TypeScript

## 📦 Installation

```bash
Delete file package-lock.json
```

```bash
# Install dependencies
npm install
```

## 🚀 Development

To start the development server:

```bash
npm run dev
```

## 🏗️ Build

To build the project for production:

```bash
npm run build
```

## 👀 Preview

Preview the production build locally:

```bash
npm run preview
```

## 🧹 Linting

To run ESLint on the codebase:

```bash
npm run lint
```

## 📚 Sample Data

Here is example data used in the application:
- **Location File**: src/constans/form

```ts
export const dataQueue = [
  { email: "john.doe@example.com", url: "http://www.google.com" },
  { email: "alice99@mail.com", url: "http://www.google.com" },
  { email: "test.user123@gmail.com", url: "http://www.google.com" },
  { email: "randomguy@hotmail.com", url: "http://www.google.com" },
  { email: "hello.world@domain.co", url: "http://www.google.com" },
  { email: "foo.bar@testing.org", url: "http://www.google.com" },
  { email: "janedoe987@yahoo.com", url: "http://www.google.com" },
  { email: "contact@newmail.com", url: "http://www.google.com" },
  { email: "support@service.tech", url: "http://www.googlessss.com" },
  { email: "user1234@nowhere.net", url: "http://www.googless.com" }
];
```

## 🧪 Testing

> Not yet implemented – consider integrating [Jest](https://jestjs.io/) or [Vitest](https://vitest.dev/).

## 📄 License

Private project – All rights reserved.

---

_Developed with ❤️ by your team._