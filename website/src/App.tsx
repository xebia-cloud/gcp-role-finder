import {
  Admin,
  Resource,
  ListGuesser,
  EditGuesser,
  ShowGuesser,
} from "react-admin";
import { RoleList, RoleShow } from "./Roles";
import { dataProvider } from "./dataProvider";

export const App = () => (
  <Admin dataProvider={dataProvider}>
    <Resource
      name="roles"
      list={RoleList}
      show={RoleShow}
    />
  </Admin>
);

