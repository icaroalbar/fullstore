import { MainLayout } from "@/layouts/MainLayout"
import { Catalog } from "@/features/catalog/Catalog"

export function App() {
  return (
    <MainLayout>
      <Catalog />
    </MainLayout>
  )
}

export default App
