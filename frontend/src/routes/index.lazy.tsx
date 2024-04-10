import { FormatSelect } from "@/components/format-select";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { createLazyFileRoute } from "@tanstack/react-router";
import download from "downloadjs";
import { useState } from "react";

export const Route = createLazyFileRoute("/")({
  component: Index,
});

type Format = "audio" | "video";

function Index() {
  const [url, setUrl] = useState("");
  const [format, setFormat] = useState<Format>();

  const [isFormError, setFormError] = useState(false);

  const submitDownloadRequest = () => {
    if (!url || !format) {
      setFormError(true);
      return;
    }

    setFormError(false);

    fetch("http://localhost:8080/download", {
      method: "POST",
      mode: "cors",
      body: JSON.stringify({
        url: url,
        format: format,
      }),
    })
      .then((res) => res.blob())
      .then((blob) => {
        download(blob);
      });
  };

  return (
    <div className="">
      <h1 className="text-xl my-2">Youtube downloader</h1>
      <div className="flex flex-col gap-2">
        <Input
          placeholder="Youtube URL"
          onInput={(e: React.ChangeEvent<HTMLInputElement>) =>
            setUrl(e.target.value)
          }
        />
        <FormatSelect onValueChange={(value) => setFormat(value as Format)} />
        {isFormError && (
          <div className="px-1 text-red-500">
            <span>Give a url and set the format.</span>
          </div>
        )}
        <Button className="self-start" onClick={submitDownloadRequest}>
          Download
        </Button>
      </div>
    </div>
  );
}
