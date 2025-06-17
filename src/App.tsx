import React, { useEffect } from "react";
import { useDispatch } from "react-redux";
import {
  Container,
  Box,
} from "@mui/material";
import { AppDispatch } from "./store/store";
import { Navigate, Route, Routes } from "react-router-dom";
import QueueList from "./page/queueList";

const App: React.FC = () => {
  const dispatch = useDispatch<AppDispatch>();

  useEffect(() => {
    const storedUser = localStorage.getItem("user");
    if (storedUser) {
      const parsedUser = JSON.parse(storedUser);
      if (parsedUser?.accessToken) {
        dispatch({
          type: "auth/loginAsync/fulfilled",
          payload: parsedUser,
        });
      }
    }
  }, [dispatch]);

  return (
    <>
     <Container maxWidth="xl">
        <Box sx={{ my: 4 }}>
           <Routes>
            <Route path="/" element={<Navigate to="/forms" />} />
            <Route path="/forms" element={<QueueList />} />
          </Routes>
        </Box>
      </Container>
    </>
  );
};

export default App;
