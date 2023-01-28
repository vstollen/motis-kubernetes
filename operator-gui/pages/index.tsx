import { useState } from "react";
import FormCheckboxInput from "../components/FormCheckboxInput";
import FormCodeInput from "../components/FormCodeInput";
import FormTextInput from "../components/FormTextInput";

export default function Home() {
  const [name, setName] = useState("");
  const [scheduleUrl, setScheduleUrl] = useState("");
  const [osmUrl, setOsmUrl] = useState("");
  const [periodicRefreshs, setPeriodicRefreshs] = useState(false);
  const [refreshCron, setRefreshCron] = useState("");
  const [configIni, setConfigIni] = useState("");

  function handleSubmit() {
    fetch("/api/motis", {
      method: "POST",
      headers: {
        "Accept": "application/json",
        "Content-Type": "application/json"
      },
      body: JSON.stringify({
        name: name,
        scheduleUrl: scheduleUrl,
        osmUrl: osmUrl,
        refreshSchedule: periodicRefreshs ? refreshCron : null,
        config: configIni
      })
    });
  }

  return (
    <div className="flex justify-center p-20">
      <section className="max-w-2xl flex-auto rounded-md bg-zinc-50 p-16">
        <div className="">
          <h1 className="mb-16 text-xl font-bold">Create MOTIS instance</h1>
          <div className="flex flex-col gap-8">
            <FormTextInput
              type="text"
              id="name"
              value={name}
              onChange={setName}
              label="Name:"
              placeholder="motis-hessen"
            />
            <FormTextInput
              type="url"
              id="schedule"
              value={scheduleUrl}
              onChange={setScheduleUrl}
              label="Schedule URL:"
              placeholder="https://opendata.avv.de/current_GTFS/AVV_GTFS_mit_SPNV.zip"
            />
            <FormTextInput
              type="url"
              id="osm"
              value={osmUrl}
              onChange={setOsmUrl}
              label="Open Street Map (OSM):"
              placeholder="https://download.geofabrik.de/europe/germany/nordrhein-westfalen/koeln-regbez-latest.osm.pbf"
            />
          </div>
          <hr className="my-8" />
          <div className="flex flex-col">
            <FormCheckboxInput
              id="periodicRefreshs"
              label="Periodic Refreshs"
              checked={periodicRefreshs}
              onChange={() => setPeriodicRefreshs(!periodicRefreshs)}
            />
            <div className={periodicRefreshs ? "" : "hidden"}>
              <FormTextInput
                type="text"
                id="refreshCron"
                value={refreshCron}
                onChange={setRefreshCron}
                placeholder="0 3 * * *"
              />
            </div>
          </div>
          <hr className="my-8" />
          <FormCodeInput id="config" label="config.ini:" value={configIni} onChange={setConfigIni} />
          <button
            className="float-right mt-8 rounded-md bg-rose-600 p-4 py-3 font-bold text-[hsl(347_100%_99%)] shadow shadow-rose-600/50 transition-all hover:bg-rose-500 hover:shadow-md hover:shadow-rose-500/50 active:shadow"
            onClick={handleSubmit}
          >
            Create MOTIS instance
          </button>
        </div>
      </section>
    </div>
  );
}
