import React from "react";
import {
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  Paper,
  Typography,
  CircularProgress,
  Alert,
  Button,
  Box,
  Chip,
  IconButton,
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  DialogContentText,
  TextField,
} from "@mui/material";
import ReplayIcon from "@mui/icons-material/Replay";
import moment from "moment";
import { useQuery } from "@tanstack/react-query";
import { useDispatch } from "react-redux";
import { AppDispatch } from "../store/store";
import { createQueueAsync, deleteQueueAsync, getQueueAsync, retryQueueAsync } from "../store/queueSlice";
import CloseIcon from '@mui/icons-material/Close';
import { dataQueue } from "../constants/form";

const QueueList: React.FC = () => {
  const dispatch = useDispatch<AppDispatch>();
   const [retryConfirmOpen, setRetryConfirmOpen] = React.useState(false);
   const [questionToRetry, setQuestionToRetry] = React.useState<any>(null);
    const [delayMs, setDelayMs] = React.useState<number>(1000);
   const delay = (ms: number) => new Promise(resolve => setTimeout(resolve, ms));

  const {
    data: queue,
    isLoading,
    isError,
    error,
  } = useQuery({
    queryKey: ["queue"],
    queryFn: async () => {
      return await dispatch(getQueueAsync()).unwrap();
    },
    refetchInterval: 1500,
  });

  const confirmRetry = () => {
    dispatch(retryQueueAsync({ formData: questionToRetry }));

    setRetryConfirmOpen(false);

  };

const addQueue = async () => {
  for (let i = 0; i < dataQueue.length; i++) {
    const data = dataQueue[i];

    try {
      await dispatch(createQueueAsync({ email: data.email, url: data.url })).unwrap();
    } catch (error) {
      console.error("Failed to create queue for", data.email, "Reason:", error);
    }

    await delay(delayMs);
  }
};

  const deleteQueue = () => {
    dispatch(deleteQueueAsync());
  };

  return (
    <>
      <Box display="flex" justifyContent="space-between" mt={2}>
        <Typography variant="h6" sx={{ m: 1 }}>
          Queue List 
        </Typography>
        <Box>
          <Typography variant="h6" sx={{ m: 1 }}>
            maximum process 3 (thread)
          </Typography>
          <TextField
            label="Interval (ms) Add Queue"
            type="number"
            value={delayMs}
            onChange={(e) => setDelayMs(Number(e.target.value))}
            size="small"
            sx={{ width: 170, mr: 2, mt: 1 }}
            inputProps={{ min: 0 }}
          />
          <Button
            variant="contained"
            color="error"
            sx={{ m: 1 }}
            onClick={deleteQueue}
            disabled={queue?.length === 0}
          >
            {isLoading ? <CircularProgress size={24} color="inherit" /> : "Delete Queue"}
          </Button>
          <Button
            variant="contained"
            sx={{ m: 1 }}
            onClick={addQueue}
            disabled={queue!.length > 1 }
          >
            {isLoading ? <CircularProgress size={24} color="inherit" /> : "Add Queue"}
          </Button>
        </Box>
  
      </Box>

      <TableContainer component={Paper} sx={{ mt: 4 }}>
        {isLoading && <CircularProgress sx={{ m: 2 }} />}
        {isError && <Alert severity="error">{(error as Error).message}</Alert>}
          <Table>
            <TableHead>
              <TableRow>
                <TableCell>Email</TableCell>
                <TableCell>URL</TableCell>
                <TableCell>Status</TableCell>
                <TableCell>Created At</TableCell>
                <TableCell>Updated At</TableCell>
              </TableRow>
            </TableHead>
            <TableBody>
              {(queue || [])
                .slice()
                .sort((a: any, b: any) => new Date(a.created_at).getTime() - new Date(b.created_at).getTime())
                .map((form: any) => (
                  <TableRow key={form.id}>
                    <TableCell>{form.email}</TableCell>
                    <TableCell>{form.url}</TableCell>
                   <TableCell>
                    <Box display="flex" alignItems="center" gap={1}>
                      <Chip
                        label={form.status}
                        color={
                          form.status === "success"
                            ? "success"
                            : form.status === "failed"
                            ? "error"
                            : form.status === "pending"
                            ? "warning"
                            : "default"
                        }
                        variant="outlined"
                      />

                      {form.status === "failed" && (
                        <IconButton
                          color="primary"
                          size="small"
                          onClick={() => {
                            setQuestionToRetry(form)
                            setRetryConfirmOpen(true)
                          }}
                        >
                          <ReplayIcon fontSize="small" />
                        </IconButton>
                      )}
                    </Box>
                  </TableCell>
                    <TableCell>{moment(form.created_at).format("DD MMM YYYY, HH:mm")}</TableCell>
                    <TableCell>{moment(form.updated_at).format("DD MMM YYYY, HH:mm")}</TableCell>
                  </TableRow>
                ))}
            </TableBody>
          </Table>

        {!isLoading && queue?.length === 0 && (
          <Typography sx={{ m: 2 }}>No queues found.</Typography>
        )}
      </TableContainer>


         <Dialog open={retryConfirmOpen} onClose={() => setRetryConfirmOpen(false)}>
               <DialogTitle sx={{ m: 0, p: 2 }}>
                  Retry
                  <IconButton
                    aria-label="close"
                    onClick={() => setRetryConfirmOpen(false)}
                    sx={{
                      position: 'absolute',
                      right: 8,
                      top: 8,
                      color: (theme) => theme.palette.grey[500],
                    }}
                  >
                    <CloseIcon />
                  </IconButton>
                </DialogTitle>
              <DialogContent>
                <DialogContentText>
                  Apalah ingin mengulangi proses queue? 
                  proses ini mesimulasikan perubahan status dengan cara mengubah url menjadi google.com agar ketika ping tidak timeout sehingga status menjadi success
                </DialogContentText>
              </DialogContent>
              <DialogActions>
                <Button onClick={() => setRetryConfirmOpen(false)}>Cancel</Button>
                <Button onClick={confirmRetry} variant="contained" color="error">
                  Retry
                </Button>
              </DialogActions>
            </Dialog>
    </>
  );
};

export default QueueList;
