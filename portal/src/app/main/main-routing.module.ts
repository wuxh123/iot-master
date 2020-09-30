import {NgModule} from '@angular/core';
import {Routes, RouterModule} from '@angular/router';
import {DashComponent} from './dash/dash.component';
import {MainComponent} from './main.component';
import {LinkComponent} from './link/link.component';
import {ChannelComponent} from './channel/channel.component';
import {PluginComponent} from './plugin/plugin.component';
import {LinkMonitorComponent} from './link-monitor/link-monitor.component';
import {ModelComponent} from './model/model.component';
import {TunnelComponent} from './tunnel/tunnel.component';
import {VariableComponent} from './variable/variable.component';
import {BatchComponent} from './batch/batch.component';
import {JobComponent} from './job/job.component';
import {StrategyComponent} from './strategy/strategy.component';
import {ModelEditComponent} from "./model-edit/model-edit.component";

const routes: Routes = [
  {
    path: '',
    component: MainComponent,
    children: [
      {path: '', redirectTo: 'dash'},
      {path: 'dash', component: DashComponent},
      {path: 'channel', component: ChannelComponent},
      {path: 'link', component: LinkComponent},
      {path: 'link-monitor/:id', component: LinkMonitorComponent},
      {path: 'plugin', component: PluginComponent},
      {path: 'model', component: ModelComponent},
      {path: 'model-create', component: ModelEditComponent},
      {path: 'model-edit/:id', component: ModelEditComponent},
      {path: 'tunnel', component: TunnelComponent},
      {path: 'variable', component: VariableComponent},
      {path: 'batch', component: BatchComponent},
      {path: 'job', component: JobComponent},
      {path: 'strategy', component: StrategyComponent},
      {path: '**', redirectTo: 'dash'},
    ]
  },
];

@NgModule({
  imports: [RouterModule.forChild(routes)],
  exports: [RouterModule]
})
export class MainRoutingModule {
}
