import { MainLayout } from "@/layouts/MainLayout"

export function App() {
  return (
    <MainLayout>
      <div className="py-12 px-4">
        <h2 className="text-3xl font-bold mb-2">Catálogo em breve</h2>
        <p className="text-muted-foreground">
          Os produtos aparecerão aqui na próxima atualização.
        </p>
      </div>
    </MainLayout>
  )
}

export default App
