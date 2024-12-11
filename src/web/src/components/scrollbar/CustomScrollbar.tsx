import React from "react";
import { Scrollbars } from "react-custom-scrollbars-2";

export const CustomScrollbar = ({ children, style }: { children: React.ReactNode; style?: React.CSSProperties }) => {
  const renderTrack = ({ style, ...props }: any) => {
    const trackStyle = {
      position: "absolute",
      maxWidth: "100%",
      width: 6,
      transition: "opacity 200ms ease 0s",
      opacity: 0,
      background: "transparent",
      bottom: 2,
      top: 2,
      borderRadius: 3,
      right: 0,
    };
    return <div style={{ ...style, ...trackStyle }} {...props} />;
  };

  const renderThumb = ({ style, ...props }: any) => {
    const thumbStyle = {
      borderRadius: 15,
      background: "rgba(222, 222, 222, .1)",
    };
    return <div style={{ ...style, ...thumbStyle }} {...props} />;
  };

  const renderView = ({ style, ...props }: any) => {
    const viewStyle = {
      marginBottom: -22,
    };
    return <div style={{ ...style, ...viewStyle }} {...props} />;
  };

  return (
    <Scrollbars
      style={style}
      renderTrackVertical={renderTrack}
      renderThumbVertical={renderThumb}
      renderView={renderView}
    >
      {children}
    </Scrollbars>
  );
};
