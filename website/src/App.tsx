import {
    Admin,
    Resource,
    ListGuesser,
    EditGuesser,
    ShowGuesser, defaultTheme,
} from "react-admin";
import { RoleList, RoleShow } from "./Roles";
import { dataProvider } from "./dataProvider";
import {Dashboard} from "./Dashboard";
import purple from '@mui/material/colors/purple';

const theme = {
    ...defaultTheme,
    palette: {
        primary: purple
    },
};

export const App = () => (
  <Admin dataProvider={dataProvider} dashboard={Dashboard} theme={theme}>
    <Resource
      name="roles"
      list={RoleList}
      show={RoleShow}
    />
  </Admin>
);

