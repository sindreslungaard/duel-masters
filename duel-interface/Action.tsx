import { useEffect, useRef } from "react";
import { Popup } from "./Popup";

export interface ActionProps {
  title: string;
  visible: boolean;
}

export function Action({ title, visible }: ActionProps) {
  return (
    <>
      <Popup title={title} visible={visible} maxWidth="500px" zIndex={1000}>
        <div className="p-6">Your content here</div>
      </Popup>
    </>
  );
}
