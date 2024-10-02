import {
  createSignal,
  createResource,
  Switch,
  Match,
  Show,
  For,
} from "solid-js";

type Record = {
  id: number;
  name: string;
  created_at: string;
};

function App() {
  const [name, setName] = createSignal("");
  const [loading, setLoading] = createSignal(false);
  const [delLoading, setDelLoading] = createSignal(0);
  const [isEditing, setIsEditing] = createSignal(0);
  const [editLoading, setEditLoading] = createSignal(false);
  const [editName, setEditName] = createSignal("");
  const [error, setError] = createSignal("");

  const [data, { mutate }] = createResource(async () => {
    const response = await fetch("http://localhost:8081/api/records");
    const json = await response.json();
    return json;
  });

  async function createRecord(name: string) {
    setLoading(true);
    setError("");
    const payload = { name: name };
    try {
      const res = await fetch("http://localhost:8081/api/records", {
        method: "POST",
        body: JSON.stringify(payload),
        headers: {
          "Content-Type": "application/json",
        },
      });
      if (res.status !== 201) {
        const errorText = await res.text();
        throw new Error(errorText);
      }
      const addedItem = await res.json();
      mutate((prev) => (prev ? [...prev, addedItem] : [addedItem]));
    } catch (err) {
      if (err instanceof Error) {
        setError(err.message);
      } else {
        setError("Something went wrong");
      }
    } finally {
      setName("")
      setLoading(false);
    }
  }

  async function editRecord(id: number, name: string) {
    setEditLoading(true);
    setError("");
    const payload = { id: id, name: name };
    try {
      const res = await fetch("http://localhost:8081/api/records", {
        method: "PUT",
        body: JSON.stringify(payload),
        headers: {
          "Content-Type": "application/json",
        },
      });
      if (res.status !== 200) {
        const errorText = await res.text();
        throw new Error(errorText);
      }
      const updatedItem = await res.json();
      mutate((prev: Record[]) =>
        prev
          ? prev.map((item) => (item.id === id ? updatedItem : item))
          : [updatedItem]
      );
    } catch (err) {
      if (err instanceof Error) {
        setError(err.message);
      } else {
        setError("Something went wrong");
      }
    } finally {
      setEditName("")
      setIsEditing(0);
      setEditLoading(false);
    }
  }

  async function deleteRecord(id: number) {
    setDelLoading(id);
    setError("");
    try {
      const res = await fetch(`http://localhost:8081/api/records/${id}`, {
        method: "DELETE",
        headers: {
          "Content-Type": "application/json",
        },
      });
      if (res.status != 204) {
        const errorText = await res.text();
        throw new Error(errorText);
      }
      mutate((prev: Record[]) =>
        prev ? prev.filter((item) => item.id !== id) : []
      );
    } catch (err) {
      if (err instanceof Error) {
        setError(err.message);
      } else {
        setError("Something went wrong");
      }
    } finally {
      setDelLoading(0);
    }
  }

  return (
    <div class="flex justify-center p-4">
      <div>
        <div class="flex gap-2">
          <input
            class="bg-zinc-900 rounded-[8px] border border-zinc-700 px-4 py-2"
            placeholder="Record name"
            value={name()}
            onInput={(e) => setName(e.currentTarget.value)}
          />
          {error() && <p>{error()}</p>}
          <button
            class="bg-indigo-500 text-white rounded-[8px] px-4 py-2"
            onClick={() => {
              createRecord(name());
            }}
          >
            {loading() ? "Loading..." : "Create"}
          </button>
        </div>

        <div class="flex flex-col gap-4">
          <Show when={data.loading}>
            <p class="text-center mt-[10px]">Loading...</p>
          </Show>
          <Switch>
            <Match when={data.error}>
              <p class="text-center mt-[10px]">Error: {data.error}</p>
            </Match>
            <Match when={data()}>
              <For each={data()}>
                {(record: Record) => (
                  <>
                    {isEditing() === record.id ? (
                      <div class="flex gap-2 p-4 border border-zinc-700 rounded-[8px] mt-[10px]">
                        <input
                          class="bg-zinc-900 rounded-[8px] border border-zinc-700 px-4 py-2"
                          placeholder="Record name"
                          value={editName()}
                          onInput={(e) => setEditName(e.currentTarget.value)}
                        />
                        <button
                          class="bg-indigo-500 text-white rounded-[8px] px-4 py-2"
                          onClick={() => {
                            editRecord(record.id, editName());
                          }}
                        >
                          {editLoading() ? "Saving..." : "Save"}
                        </button>
                      </div>
                    ) : (
                      <div class="p-4 border border-zinc-700 rounded-[8px] mt-[10px]">
                        <p>{record.id}</p>
                        <p>{record.name}</p>
                        <p>{record.created_at}</p>
                        <div class="flex gap-2">
                          <button
                            class="bg-red-500 text-white rounded-[8px] px-2 py-1"
                            onClick={() => deleteRecord(record.id)}
                          >
                            {delLoading() === record.id
                              ? "Deleting..."
                              : "Delete"}
                          </button>
                          <button
                            class="bg-indigo-500 text-white rounded-[8px] px-2 py-1"
                            onClick={() => {
                              setEditName(record.name);
                              setIsEditing(record.id);
                            }}
                          >
                            Edit record
                          </button>
                        </div>
                      </div>
                    )}
                  </>
                )}
              </For>
            </Match>
          </Switch>
        </div>
      </div>
    </div>
  );
}
export default App;
