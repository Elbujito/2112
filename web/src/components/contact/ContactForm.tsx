import React, { useState } from "react";
import axios from "axios";
import { Grid, TextField, Button, Typography, CircularProgress, Alert } from "@mui/material";

const ContactForm = () => {
  const [status, setStatus] = useState({
    submitted: false,
    submitting: false,
    info: { error: false, msg: "" },
  });
  const [inputs, setInputs] = useState({
    email: "",
    message: "",
  });

  const handleServerResponse = (ok: boolean, msg: string) => {
    if (ok) {
      setStatus({
        submitted: true,
        submitting: false,
        info: { error: false, msg: msg },
      });
      setInputs({
        email: "",
        message: "",
      });
    } else {
      setStatus({
        submitted: true,
        submitting: false,
        info: { error: true, msg: msg },
      });
    }
  };

  const handleOnChange = (e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement>) => {
    const { id, value } = e.target;
    setInputs((prev) => ({
      ...prev,
      [id]: value,
    }));
    setStatus({
      submitted: false,
      submitting: false,
      info: { error: false, msg: "" },
    });
  };

  const handleOnSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setStatus((prevStatus) => ({ ...prevStatus, submitting: true }));
    try {
      await axios.post("https://formspree.io/f/mgebnnkr", inputs);
      handleServerResponse(true, "Thank you! Your message has been submitted.");
    } catch (error: any) {
      handleServerResponse(false, error.response?.data?.error || "An error occurred.");
    }
  };

  return (
    <div
      style={{
        background: "linear-gradient(to bottom right, #1c1c1c, #2c2c2c)",
        height: "100vh",
        width: "100vw",
        display: "flex",
        justifyContent: "center",
        alignItems: "center",
        color: "#fff",
        padding: "1rem",
      }}
    >
      <div
        style={{
          background: "#282828",
          borderRadius: "8px",
          padding: "2rem",
          maxWidth: "600px",
          width: "100%",
          boxShadow: "0px 4px 10px rgba(0, 0, 0, 0.5)",
        }}
      >
        <Typography variant="h4" component="h1" style={{ textAlign: "center", marginBottom: "1.5rem", color: "#e0e0e0" }}>
          Contact Us
        </Typography>
        <form onSubmit={handleOnSubmit} noValidate>
          <Grid container spacing={3}>
            <Grid item xs={12}>
              <TextField
                variant="outlined"
                onChange={handleOnChange}
                value={inputs.email}
                required
                fullWidth
                id="email"
                label="Email"
                type="email"
                InputProps={{
                  style: { background: "#3c3c3c", borderRadius: "4px", color: "#fff" },
                }}
                InputLabelProps={{
                  style: { color: "#aaa" },
                }}
              />
            </Grid>
            <Grid item xs={12}>
              <TextField
                variant="outlined"
                onChange={handleOnChange}
                value={inputs.message}
                required
                fullWidth
                id="message"
                label="Message"
                multiline
                rows={6}
                InputProps={{
                  style: { background: "#3c3c3c", borderRadius: "4px", color: "#fff" },
                }}
                InputLabelProps={{
                  style: { color: "#aaa" },
                }}
              />
            </Grid>
          </Grid>
          <Button
            type="submit"
            fullWidth
            disabled={status.submitting}
            variant="contained"
            style={{
              marginTop: "1.5rem",
              background: "#4a90e2",
              color: "#fff",
              padding: "0.75rem",
              borderRadius: "4px",
              fontWeight: "bold",
              fontSize: "1rem",
              textTransform: "none",
            }}
          >
            {status.submitting ? <CircularProgress size={24} color="inherit" /> : "Send Message"}
          </Button>
        </form>

        {status.info.error && (
          <Alert severity="error" style={{ marginTop: "1rem", background: "#5a1c1c", color: "#ffdddd" }}>
            {status.info.msg}
          </Alert>
        )}
        {!status.info.error && status.info.msg && (
          <Alert severity="success" style={{ marginTop: "1rem", background: "#1c5a1c", color: "#ddffdd" }}>
            {status.info.msg}
          </Alert>
        )}
      </div>
    </div>
  );
};

export { ContactForm };
