import React, { useState } from 'react';
import axios from 'axios';
import Grid from '@mui/material/Grid';
import Button from '@mui/material/Button';
import TextField from '@mui/material/TextField';

const ContactForm =  () => {
  const [status, setStatus] = useState({
    submitted: false,
    submitting: false,
    info: { error: false, msg: "" },
  });
  const [inputs, setInputs] = useState({
    email: '',
    message: '',
  });
  const handleServerResponse = (ok: boolean, msg: string) => {
    if (ok) {
      setStatus({
        submitted: true,
        submitting: false,
        info: { error: false, msg: msg },
      });
      setInputs({
        email: '',
        message: '',
      });
    } 
    else {
      setStatus({
        submitted: true,
        submitting: false,
        info: { error: true, msg: msg },
      });
    }
    
  };
  const handleOnChange = (e: { persist: () => void; target: { id: any; value: any; }; }) => {
    e.persist();
    setInputs((prev) => ({
      ...prev,
      [e.target.id]: e.target.value,
    }));
    setStatus({
      submitted: false,
      submitting: false,
      info: { error: false, msg: "" },
    });
  };
  const handleOnSubmit = (e: { preventDefault: () => void; }) => {
    e.preventDefault();
    setStatus((prevStatus) => ({ ...prevStatus, submitting: true }));
    axios({
      method: 'POST',
      url: 'https://formspree.io/f/mgebnnkr',
      data: inputs,
    })
      .then(() => {
        handleServerResponse(
          true,
          'Thank you, your message has been submitted.',
        );
      })
      .catch((error) => {
        handleServerResponse(false, error.response.data.error);
      });
  };
  return (
    <div>
      <form onSubmit={handleOnSubmit} noValidate>
        <Grid container spacing={2}>
             <Grid item xs={12}>
              <TextField 
                name="_replyto"
                 variant="outlined"
                 onChange={handleOnChange}
                 value={inputs.email}
                required
                 fullWidth
                 id="email"
                 label="Email"
                 autoFocus
               />
             </Grid>
             <Grid item xs={12}>
              <textarea rows={10} className="w-full p-4 border text-black-100 border-gray-600"
                id="message"
                name="message"
                onChange={handleOnChange}
                required
                value={inputs.message}
              />
             </Grid>
        </Grid>
        <Button
             type="submit"
            fullWidth
            disabled={status.submitting}
             variant="contained"
             className="text-white-100 bg-black-100"
           >
             {!status.submitting
            ? !status.submitted
              ? 'Submit'
              : 'Submitted'
            : 'Submitting...'}
        </Button>
      </form>
      {status.info.error && (
        <div className="text-red-100 ">Error: {status.info.msg} - Please check your email address</div>
      )}
      {!status.info.error && status.info.msg && <p className="text-black-100 ">{status.info.msg}</p>}
    </div>
  );
};

export { ContactForm };